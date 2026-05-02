package player

import (
	"codeberg.org/Asep5K/animein/api"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func PlayVideo(url string, title string) {
	if os.Getenv("ANDROID_DATA") == "/data" {
		cmd := exec.Command("termux-open", url)
		if err := cmd.Start(); err != nil {
			fmt.Printf("✘ Gagal buka Termux: %v\n", err)
			time.Sleep(5 * time.Second)
		}
		return
	}

	args := []string{
		url, "--referrer=" + api.BaseURL,
		"--cache=yes", "--title=Animein CLI - " + title,
		"--force-media-title=" + title,
		"--save-watch-history=yes",
		"--save-position-on-quit=yes",
	}

	cmd := exec.Command("mpv", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("✘ Gagal muter MPV: %v\n", err)
		time.Sleep(5 * time.Second)
	}
}

// vim: ft=go
