package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/codenoid/pc-auth-notifier/shared-packages/model"
	"github.com/gin-gonic/gin"
	"github.com/tbalthazar/onesignal-go"
	"gorm.io/gorm"
)

var (
	bindAddr string
	mainDB   *gorm.DB

	osAppID = os.Getenv("OS_APP_ID") // Application -> Settings -> Keys & IDs ->ONESIGNAL APP ID
)

func main() {

	flag.StringVar(&bindAddr, "bind", ":8080", "-bind 127.0.0.1:8080")
	flag.Parse()

	pushNotif := onesignal.NewClient(nil)
	pushNotif.AppKey = os.Getenv("OS_APP_SECRET") // Application -> Settings -> Keys & IDs -> REST API KEY
	pushNotif.UserKey = os.Getenv("OS_USER_KEY")  // User Auth Key from Profile Picture -> Account & API Keys

	r := gin.Default()

	notif := r.Group("/notification")
	{
		notif.GET("/list", func(c *gin.Context) {
			machineID := c.Query("id")
			authLog := make([]model.AuthLog, 0)
			where := mainDB.Where("machine_id = ? order by timestamp desc limit 20", machineID)
			if tx := where.Find(&authLog); tx.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"message": tx.Error.Error(),
					"logs":    make([]string, 0),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"error":   false,
				"message": "",
				"logs":    authLog,
			})
		})
		notif.POST("/add", func(c *gin.Context) {
			log := model.AuthLog{}
			if err := c.BindJSON(&log); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   true,
					"message": err.Error(),
				})
				return
			}
			if tx := mainDB.Create(&log); tx.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"message": tx.Error.Error(),
				})
				return
			}
			pushNotif.Notifications.Create(&onesignal.NotificationRequest{
				AppID:            osAppID,
				Headings:         map[string]string{"en": "Someone tried to unlock your device"},
				Contents:         map[string]string{"en": log.Raw},
				IncludedSegments: []string{"Subscribed Users"},
				BigPicture:       "https://pbs.twimg.com/media/EUdhmkOXgAAemsu.jpg",
			})
			c.JSON(http.StatusOK, gin.H{
				"error":   false,
				"message": "",
			})
			return
		})
	}
	r.Run(bindAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
