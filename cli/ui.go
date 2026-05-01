package cli

import (
	"animein/api"
	"animein/models"
	"animein/player"
	"animein/utils"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/manifoldco/promptui"
)

func StartApp() {
	var query string
	if len(os.Args) > 1 {
		query = strings.Join(os.Args[1:], " ")
	}

	animeID, title := GetAnimeID(query)
	if animeID == "" {
		return
	}

	count, err := api.GetPageCount(animeID)
	if err != nil {
		return
	}
	ProcessAndSelect(animeID, title, count)
}

func GetAnimeID(initialQuery string) (string, string) {

	var result []models.Movie
	if initialQuery != "" {
		res, err := api.SearchAnime(initialQuery)
		if err != nil {
			fmt.Printf("✘ %v\n", err)
			result = TrySearch()
		} else {
			result = res
		}
	} else {
		result = TrySearch()
	}

	if len(result) == 0 {
		return "", ""
	}

	var titles []string
	for _, anime := range result {
		titles = append(titles, anime.Title)
	}

	// settup promp
	prompt := promptui.Select{
		Label: "Pilih anime",
		Items: titles,
		Size:  10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("✘ Exit: %v\n", err)
		return "", ""
	}

	return result[idx].ID, result[idx].Title
}

func TrySearch() []models.Movie {
	for i := 0; i < 3; i++ {
		input, err := utils.InputUser("Masukan judul")
		if err != nil {
			return nil
		}
		utils.ClearScreen()
		res, err := api.SearchAnime(input)
		if err == nil {
			return res
		}
		fmt.Printf("✘ %v\n", err)
	}
	return nil
}

func ProcessAndSelect(animeID string, animeTitle string, pageCount int) {
	stopLoading := utils.Loading("Fetching episodes for " + animeTitle)
	var wg sync.WaitGroup
	var mu sync.Mutex

	allResults := make(map[string]models.FinalData)
	var idList []string
	var epLabels []string

	for i := pageCount; i >= 0; i-- {
		episodesPage, err := api.GetEpisodesCached(animeID, i)
		if err != nil {
			fmt.Printf("\033[31m[!]\033[0m Skip halaman %d gara-gara error: %v\n", i, err)
			continue
		}
		results := api.ParseEpisodes(episodesPage)
		for ep := range results {
			wg.Add(1)
			mu.Lock()
			idList = append(idList, ep.ID)
			epLabels = append(epLabels, fmt.Sprintf("%s %s", animeTitle, ep.EpTitle))
			mu.Unlock()

			go func(id string, idx string) {
				defer wg.Done()
				info := api.GetEpisodeInfo(id)
				mu.Lock()
				allResults[id] = models.FinalData{
					Info:    info,
					EpTitle: idx,
				}
				mu.Unlock()
			}(ep.ID, ep.EpTitle)
		}
	}
	wg.Wait()
	stopLoading <- true

	epLabels = append(epLabels, "Keluar") // Tambahin opsi keluar

	for { // infinity loop
		utils.ClearScreen()
		prompt := promptui.Select{
			Label: "Pilih Episode",
			Items: epLabels,
			Size:  15,
		}

		idx, epInfo, err := prompt.Run()
		if err != nil {
			return
		}
		// Kalau milih opsi paling bawah (Keluar)
		if epInfo == "Keluar" {
			return
		}

		selectedID := idList[idx]
		dataEpisode := allResults[selectedID]

		var directServers []models.Server
		for _, s := range dataEpisode.Info {
			if s.Type == "direct" {
				directServers = append(directServers, s)
			}
		}

		if len(directServers) == 0 {
			fmt.Println("Gak ada link direct!")
			time.Sleep(2 * time.Second)
			continue
		}

		var resLabels []string
		for _, s := range directServers {
			resLabels = append(resLabels, s.Quality)
		}
		resLabels = append(resLabels, "Kembali")
		resPrompt := promptui.Select{
			Label: "Pilih resolusi",
			Items: resLabels,
		}

		resIdx, res, err := resPrompt.Run()
		if err != nil {
			return
		}
		if res == "Kembali" {
			continue
		}
		fmt.Println(epInfo)
		player.PlayVideo(directServers[resIdx].Link, epInfo)
	}
}

// vim: ft=go
