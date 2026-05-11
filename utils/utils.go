package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/manifoldco/promptui"
)

func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

func InputUser(Label string) (string, error) {
	prompt := promptui.Prompt{
		Label: Label,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("✘ Prompt failed: %w", err)
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

func getHistPath() string {
	var histPath string
	fileName := "animein-state.json"
	dir, err := os.UserCacheDir()
	if err != nil {
		histPath = "." + fileName
		return histPath
	}
	histPath = filepath.Join(dir, fileName)
	return histPath
}

func SaveState(aniID string, epID string, title string) {
	histPath := getHistPath()
	state := map[string]string{
		"movie_id":   aniID,
		"episode_id": epID,
		"title":      title,
	}
	data, _ := json.MarshalIndent(state, "", " ")
	_ = os.WriteFile(histPath, data, 0644)
}

func LoadState() (map[string]string, error) {
	histPath := getHistPath()
	file, err := os.ReadFile(histPath)
	if err != nil {
		return nil, err
	}
	var state map[string]string
	json.Unmarshal(file, &state)
	return state, nil
}

func ShowState() {
	showState := flag.Bool("show-state", false, "Tampilkan history tontonan terakhir")
	flag.BoolVar(showState, "l", false, "Tampilkan history (shorthand)")
	flag.Parse()
	if *showState {
		last, err := LoadState()
		if err != nil {
			fmt.Println("Belum ada history tontonan.")
			os.Exit(0)

		}
		fmt.Printf("Terakhir ditonton: %s\n", last["title"])
		fmt.Printf("Episode ID: %s\n", last["episode_id"])
		fmt.Printf("Movie ID: %s\n", last["movie_id"])
		os.Exit(0)
	}
}

func JsonDecoder[T any](res *http.Response) (T, error) {
	var data T
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("Decoder Error: %w", err)
	}
	return data, nil
}

// vim: ft=go
