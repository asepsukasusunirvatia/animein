package cli

import (
	"animein/api"
	"animein/models"
	"animein/player"
	"animein/utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
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
		// fmt.Printf("%d. %s\n", i+1, anime.Title)
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
	var reader = bufio.NewReader(os.Stdin)
	for i := 0; i < 3; i++ {
		input := utils.InputUser("Masukan judul anime: ", reader)
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
	var wg sync.WaitGroup
	var mu sync.Mutex

	allResults := make(map[string][]models.Server)
	var idList []string

	fmt.Println("Loading episodes...")
	for i := pageCount; i >= 0; i-- {
		episodesPage, err := api.GetEpisodesCached(animeID, i)
		if err != nil {
			fmt.Printf("\033[31m[!]\033[0m Skip halaman %d gara-gara error: %v\n", i, err)
			continue
		}
		epIDChan := api.ParseEpisodes(episodesPage)

		for epID := range epIDChan {
			wg.Add(1)
			mu.Lock()
			idList = append(idList, epID)
			mu.Unlock()

			go func(id string) {
				defer wg.Done()
				info := api.GetEpisodeInfo(id)
				mu.Lock()
				allResults[id] = info
				mu.Unlock()
			}(epID)
		}
		wg.Wait()
	}

	for { // infinity loop
		utils.ClearScreen()

		var epLabels []string
		for i := range idList {
			epLabels = append(epLabels, fmt.Sprintf("%s Episode %d", animeTitle, i+1))

		}
		epLabels = append(epLabels, "Keluar (Quit)") // Tambahin opsi keluar

		prompt := promptui.Select{
			Label: "Pilih Episode",
			Items: epLabels,
			Size:  15,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return
		}
		// Kalau milih opsi paling bawah (Keluar)
		if idx == len(idList) {
			return
		}

		selectedID := idList[idx]
		servers := allResults[selectedID]
		var directServers []models.Server
		for _, s := range servers {
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
		resIdx, _, err := resPrompt.Run()
		if err != nil {
			continue
		}
		if resIdx == len(directServers) {
			continue
		}
		player.PlayVideo(directServers[resIdx].Link, animeTitle, " Ep "+strconv.Itoa(idx+1))
	}
}
