package cli

import (
	"animein/api"
	"animein/models"
	"animein/player"
	"animein/utils"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

func StartApp() {
	var query string
	if len(os.Args) > 1 {
		query = strings.Join(os.Args[1:], " ")
	}

	animeID, title := getAnimeID(query)
	if animeID == "" {
		return
	}

	count, err := api.GetPageCount(animeID)
	if err != nil {
		fmt.Println(err)
		return
	}
	processAndSelect(animeID, title, count)
}

func getAnimeID(initialQuery string) (string, string) {
	var result []models.Movie
	if initialQuery != "" {
		res, err := api.SearchAnime(initialQuery)
		if err != nil {
			fmt.Printf("✘ %v\n", err)
			result = trySearch()
		} else {
			result = res
		}
	} else {
		result = trySearch()
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

func trySearch() []models.Movie {
	for i := 0; i < 3; i++ {
		input, err := utils.InputUser("\033[94mMasukan judul")
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

func getSortedEpisodes(animeID string, page int) ([]string, []string, error) {
	episodesPage, err := api.GetEpisodesCached(animeID, page)
	if err != nil {
		return nil, nil, fmt.Errorf("ERORR: %v", err)
	}

	var ids, labels []string
	episodes := api.ParseEpisodes(episodesPage)
	for ep := range episodes {
		ids = append(ids, ep.ID)
		labels = append(labels, ep.EpTitle)
	}

	slices.Reverse(ids)
	slices.Reverse(labels)
	return ids, labels, nil
}

func selectResolution(episodeID string) (string, error) {
	epsInfo := api.GetEpsInfo(episodeID)
	var labels, links []string

	for _, srv := range epsInfo {
		if srv.Type == "direct" {
			labels = append(labels, srv.Quality)
			links = append(links, srv.Link)
		}
	}
	labels = append(labels, "Kembali")

	prompt := promptui.Select{
		Label: "Pilih resolusi:",
		Items: labels,
	}
	idx, str, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("errr %w", err)
	}
	if str == "Kembali" {
		return str, nil
	}
	return links[idx], err
}

func processAndSelect(animeID string, animeTitle string, pageCount int) {
	currentPage := pageCount
	nextPage, prevPage := "Page selanjut nya >>", "Page sebelum nya <<"

	for {
		utils.ClearScreen()
		idList, epLabels, err := getSortedEpisodes(animeID, currentPage)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Tambah navigasi
		if currentPage > 0 {
			epLabels = append(epLabels, nextPage)
		}
		if currentPage < pageCount {
			epLabels = append(epLabels, prevPage)
		}
		epLabels = append(epLabels, "Keluar")
		prompt := promptui.Select{
			Label: fmt.Sprintf("Pilih Episode (Page %d):", currentPage),
			Items: epLabels,
			Size:  12,
		}

		idx, resultStr, err := prompt.Run()
		if err != nil {
			return
		}

		// Navigasi Page
		switch resultStr {
		case nextPage:
			currentPage--
			continue
		case prevPage:
			currentPage++
			continue
		case "Keluar":
			utils.ClearScreen()
			return
		}

		// Ambil url & Play
		url, err := selectResolution(idList[idx])
		fmt.Println(url)
		if err != nil {
			fmt.Printf("Error: %v", err)
			time.Sleep(5 * time.Second)
			return
		}
		if url == "Kembali" {
			continue
		}
		player.PlayVideo(url, animeTitle)
	}
}

// vim: ft=go
