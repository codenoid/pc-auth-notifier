package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

var logDatetLayout = "Jan 2 15:04:05 2006"

func startReporter() {
	t, err := tail.TailFile("/var/log/auth.log", tail.Config{Follow: true, ReOpen: true, MustExist: true})
	if err != nil {
		panic(err)
	}
	for line := range t.Lines {
		if strings.Contains(line.Text, "gdm-password]") {
			// Feb 13 13:30:26 frestea gdm-password]: gkr-pam: unlocked login keyring
			// Feb 13 13:30:22 frestea gdm-password]: pam_unix(gdm-password:auth): authentication failure; logname= uid=0 euid=0 tty=/dev/tty1 ruser= rhost=  user=ken
			fmt.Println(line.Text)

			info := strings.Split(line.Text, " gdm-password]")[0] // Feb 13 13:30:26 frestea

			dateArr := strings.Split(info, " ") // [Feb 13 13:30:26 frestea]
			dateArr = dateArr[:len(dateArr)-1]  // [Feb 13 13:30:26]

			year := time.Now().Year()                // 2021
			dateStr := strings.Join(dateArr, " ")    // "Feb 13 13:30:26"
			dateStr = fmt.Sprint(dateStr, " ", year) // "Feb 13 13:30:26 2021"

			t, err := time.Parse(logDatetLayout, dateStr)
			if err != nil {
				log.Println(err)
				continue
			}
			if t.Unix() > startTime.Unix() {
				if strings.Contains(line.Text, "unlocked login keyring") {
					log.Println("log")
				} else if strings.Contains(line.Text, "authentication failure") {
					log.Println("authentication failure")
				}
			}

		}
	}
}
