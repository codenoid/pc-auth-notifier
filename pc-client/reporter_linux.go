package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/prop"

	"github.com/codenoid/pc-auth-notifier/shared-packages/model"
	"github.com/hpcloud/tail"
)

var logDatetLayout = "Jan 2 15:04:05 -0700 2006"

func startReporter() {
	t, err := tail.TailFile("/var/log/auth.log", tail.Config{Follow: true, ReOpen: true, MustExist: true})
	if err != nil {
		panic(err)
	}
	for line := range t.Lines {
		if strings.Contains(line.Text, "gdm-password]") {
			// Feb 13 13:30:26 frestea gdm-password]: gkr-pam: unlocked login keyring
			// Feb 13 13:30:22 frestea gdm-password]: pam_unix(gdm-password:auth): authentication failure; logname= uid=0 euid=0 tty=/dev/tty1 ruser= rhost=  user=ken

			info := strings.Split(line.Text, " gdm-password]")[0] // Feb 13 13:30:26 frestea

			dateArr := strings.Split(info, " ") // [Feb 13 13:30:26 frestea]
			hostname := dateArr[len(dateArr)-1]
			dateArr = dateArr[:len(dateArr)-1] // [Feb 13 13:30:26]

			year := time.Now().Year()                      // 2021
			dateStr := strings.Join(dateArr, " ")          // "Feb 13 13:30:26"
			dateStr = fmt.Sprint(dateStr, " +0700 ", year) // "Feb 13 13:30:26 2021"

			t, err := time.Parse(logDatetLayout, dateStr)
			if err != nil {
				log.Println(err)
				continue
			}
			if t.After(startTime) {
				username := ""
				if spltUser := strings.Split(line.Text, " user="); len(spltUser) > 1 {
					username = spltUser[1]
				}
				newLog := model.AuthLog{
					LogID:     GetMD5Hash(machineID + ":" + line.Text),
					Host:      hostname,
					Username:  username,
					Status:    "",
					Timestamp: t.Unix(),
					MachineID: machineID,
					Raw:       line.Text,
				}

				if strings.Contains(line.Text, "unlocked login keyring") {
					newLog.Status = "success"
				} else if strings.Contains(line.Text, "authentication failure") {
					newLog.Status = "failure"
				}
				if newLog.Status != "" {

					mediaStream, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
						Video: func(constraint *mediadevices.MediaTrackConstraints) {
							constraint.Width = prop.Int(600)
							constraint.Height = prop.Int(400)
						},
					})
					// Since track can represent audio as well, we need to cast it to
					// *mediadevices.VideoTrack to get video specific functionalities
					track := mediaStream.GetVideoTracks()[0]
					videoTrack := track.(*mediadevices.VideoTrack)
					defer videoTrack.Close()

					// Create a new video reader to get the decoded frames. Release is used
					// to return the buffer to hold frame back to the source so that the buffer
					// can be reused for the next frames.
					videoReader := videoTrack.NewReader(false)
					frame, release, _ := videoReader.Read()
					defer release()

					// Since frame is the standard image.Image, it's compatible with Go standard
					// library. For example, capturing the first frame and store it as a jpeg image.
					// Create a Resty Client
					client := resty.New()

					b := bytes.Buffer{}
					jpeg.Encode(&b, frame, nil)

					uploadedFile := fileIO{}

					resp, err := client.R().
						SetBody(b.Bytes()).
						Post("https://file.io")
					if err != nil {
						log.Println(err)
					} else {
						json.Unmarshal(resp.Body(), &uploadedFile)
					}

					newLog.ImageURL = uploadedFile.Link

					jsonValue, _ := json.Marshal(newLog)
					_, err = http.Post(serverAddr+"/notification/add", "application/json", bytes.NewBuffer(jsonValue))
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}

		}
	}
}
