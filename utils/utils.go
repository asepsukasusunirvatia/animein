package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

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

/*
	func InputUser(prompt string, reader *bufio.Reader) string {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		return strings.TrimSpace(input)
	}
*/

func ClearScreen() {
	var cmd *exec.Cmd

	// aku gak tau ini jalan di windows apa engga
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		// Untuk Linux dan macOS
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
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

// vim: ft=go
