package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/exec"
	"runtime"
)

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func createFolderIfNotExist() {
	_ = os.Mkdir(dbPath, 0777)
}

func toggleConfig(config string) {
	createFolderIfNotExist()
	config = dbPath + config
	if _, err := os.Stat(config); err == nil {
		os.Remove(config)
	} else {
		os.OpenFile(config, os.O_RDONLY|os.O_CREATE, 0666)
	}
}

func getConfig(config string) bool {
	createFolderIfNotExist()
	config = dbPath + config
	if _, err := os.Stat(config); err == nil {
		return true
	}
	return false
}

// splitSubN from https://stackoverflow.com/a/61469854/13961710
func splitSubN(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}

// openFile opens the specified target in the default application of the user.
func openFile(target string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, target)
	return exec.Command(cmd, args...).Start()
}

// GetMD5Hash get md5 hash from string input
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
