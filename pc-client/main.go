package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/codenoid/pc-auth-notifier/pc-client/icon"
	"github.com/denisbrodbeck/machineid"
	"github.com/getlantern/systray"
	"github.com/skip2/go-qrcode"
)

// dbPath refer to database path (folder)
var dbPath = userHomeDir() + "/.auth-notifier/"

// machineID hold unique .ProtectedID
var machineID string

// serverAddr
var serverAddr = os.Getenv("PAN_HOST")

// startTime hold app start time
var startTime = time.Now()

func init() {
	if id, err := machineid.ProtectedID("auth-notifier"); err != nil {
		log.Fatal(err)
	} else {
		machineID = id
	}
}

func main() {
	createFolderIfNotExist()
	if err := qrcode.WriteFile(machineID, qrcode.Medium, 320, dbPath+"qr.png"); err != nil {
		panic(err)
	}
	go systray.Run(onReady, onExit)
	startReporter()
}

func onReady() {
	systray.SetIcon(icon.Icon)
	systray.SetTooltip("Send Notification to your mobile phone when someone fail (or success) login to this machine")

	// convert baf1216865df494880f8c719c26613c5 to [0:12] -> BA F1 21 68 65 DF
	hashStr := machineID
	hashStr = strings.ToUpper(hashStr[0:12])
	hashStr = strings.Join(splitSubN(hashStr, 2), " ")
	showQR := systray.AddMenuItem(hashStr+" [SHOW QR]", "Your [0:12] unique machine ID")

	notif := systray.AddMenuItemCheckbox("Ignore root access", "Enable/Disable", getConfig("notif_on_success"))
	ignoreRoot := systray.AddMenuItemCheckbox("Notif on success", "Enable/Disable", getConfig("ignore_root_access"))
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		// listen for click event
		select {
		case <-showQR.ClickedCh:
			openFile(dbPath + "qr.png")
		case <-notif.ClickedCh:
			toggleConfig("notif_on_success")
		case <-ignoreRoot.ClickedCh:
			toggleConfig("ignore_root_access")
		case <-mQuit.ClickedCh:
			log.Println("exiting program...")
			os.Exit(0)
		}
	}
}

func onExit() {
	// clean up here
}
