package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Convert string to integer
func StrToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Failed to converting to integer")
		return -1
	}
	return num
}

func InputUser(prompt string, reader *bufio.Reader) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

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

// vim: ft=go
