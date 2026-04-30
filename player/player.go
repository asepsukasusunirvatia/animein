package player

import (
	"animein/api"
	"animein/utils"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func PlayVideo(url string, title string, epInfo string) {
	utils.ClearScreen()
	if os.Getenv("ANDROID_DATA") == "/data" {
		cmd := exec.Command("termux-open", url)
		if err := cmd.Start(); err != nil {
			fmt.Printf("✘ Gagal buka Termux: %v\n", err)
		}
		return
	}

	argTitle := "--title=Animein CLI - " + title + epInfo
	argMediaTitle := "--force-media-title=" + title + epInfo

	cmd := exec.Command("mpv", url, "--referrer="+api.BaseURL, "--cache=yes", argMediaTitle, argTitle)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Gagal muter MPV: %v\n", err)
		time.Sleep(3 * time.Second)
	}
}

// vim: ft=go
