package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

// Convert string to integer
func StrToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("✘ Failed to converting to integer")
		return -1
	}
	return num
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func InputUser(Label string) (string, error) {
	prompt := promptui.Prompt{
		Label: Label,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("✘ Prompt failed: %v", err)
	}
	return result, nil
}

func Loading(msg string) chan bool {
	stop := make(chan bool)
	go func() {
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		for i := 0; ; i = (i + 1) % len(chars) {
			select {
			case <-stop:
				fmt.Print("\r\033[K")
				return
			default:
				fmt.Printf("\r\033[36m%s\033[0m %s...", chars[i], msg)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	return stop
}

// vim: ft=go
