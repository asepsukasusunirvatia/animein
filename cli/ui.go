package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"codeberg.org/Asep5K/animein/api"
	"codeberg.org/Asep5K/animein/models"
	"codeberg.org/Asep5K/animein/player"
	"codeberg.org/Asep5K/animein/utils"

	"github.com/manifoldco/promptui"
)

const next, prev = "Next ", "Previous "

func StartApp() {
	utils.ShowState()
	var query string
	if len(os.Args) > 1 {
		query = strings.Join(os.Args[1:], " ")
	}

	animeID, title := getAniID(query)
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

func getAniID(initialQuery string) (string, string) {
	var result []models.Movie
	if initialQuery != "" {
		res, err := api.SearchAnime(initialQuery, 0)
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

	var details []models.Movie
	for _, anime := range result {
		details = append(details, models.Movie{ID: anime.ID, Title: anime.Title, Genre: anime.Genre})
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\033[96m󰣇\033[0m {{ .Title | cyan }}",
		Inactive: "   {{ .Title | faint }}",
		Selected: "✔ {{ .Title | green }}",
		Details: `
{{ "Id:" | faint }} {{ .ID }}
{{ "Title:" | faint }} {{ .Title }}
{{ "Genre:" | faint }} {{ .Genre }}`,
	}

	// settup promp
	prompt := promptui.Select{
		Label:     "Pilih anime ",
		Items:     details,
		Size:      10,
		Searcher:  searcher(details),
		Templates: templates,
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
		res, err := api.SearchAnime(input, 0)
		if err == nil {
			return res
		}
		fmt.Printf("✘ %v\n", err)
	}
	return nil
}

func getSortedEpisodes(aniID string, title string, page int) ([]models.Movie, error) {
	stop := utils.Loading("Wait a minute")
	episodesPage, err := api.GetEpisodesCached(aniID, page)
	defer func() {
		stop <- true
	}()
	if err != nil {
		return nil, fmt.Errorf("✘ ui.getSortedEpisodes ERROR: %w", err)
	}

	episodes := api.ParseEpisodes(episodesPage)
	n := len(episodes)

	eps := make([]models.Movie, 0, n)
	for i := n - 1; i >= 0; i-- {
		ep := episodes[i]
		eps = append(eps, models.Movie{
			ID:    ep.ID,
			Title: fmt.Sprintf("%s [%s]", title, ep.EpTitle),
		})
	}
	return eps, nil
}

func parseFileSize(sizeStr string) string {
	float, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return "Unknown filesize"
	}
	return fmt.Sprintf("%.1f MB", float)
}

func selectResolution(epID string) (string, error) {
	epsInfo := api.GetEpsInfo(epID)
	var cursorPos int
	type Entry struct {
		Link  string
		Label string
		Fs    string
	}
	entries := make([]Entry, 0, len(epsInfo)+1)
	for _, srv := range epsInfo {
		if srv.Type == "direct" {
			entries = append(entries, Entry{
				Link:  srv.Link,
				Label: srv.Quality,
				Fs:    parseFileSize(srv.FileSize),
			})
		}
	}
	entries = append(entries, Entry{Link: "Kembali", Label: "Kembali", Fs: "Kembali ke menu episode"})

	prompt := promptui.Select{
		Label:     "Pilih resolusi",
		Items:     entries,
		CursorPos: cursorPos,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "▸ {{ .Label | cyan }} [{{ .Fs | blue }}]",
			Inactive: "   {{ .Label | faint }}",
			Selected: "✔ {{ .Label | green }}",
		},
	}
	idx, _, err := prompt.Run()
	cursorPos = idx
	if err != nil {
		return "", fmt.Errorf("✘ ui.selectResolution Error: %w", err)
	}

	return entries[idx].Link, nil
}

func appendEplist(cur int, count int, epl []models.Movie) []models.Movie {
	if cur < count {
		prevBtn := models.Movie{Title: prev}
		epl = append([]models.Movie{prevBtn}, epl...)
	}
	if cur > 0 {
		nextBtn := models.Movie{Title: next}
		epl = append([]models.Movie{nextBtn}, epl...)
	}
	return epl
}

func searcher(list []models.Movie) func(string, int) bool {
	return func(input string, index int) bool {
		name := strings.Replace(strings.ToLower(list[index].Title), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func processAndSelect(aniID string, aniTitle string, pageCount int) {
	currentPage := pageCount
	var cursorPos int
	for {
		utils.ClearScreen()
		epList, err := getSortedEpisodes(aniID, aniTitle, currentPage)
		if err != nil {
			fmt.Println(err)
			return
		}
		epList = appendEplist(currentPage, pageCount, epList)

		prompt := promptui.Select{
			Label:     fmt.Sprintf("Pilih Episode [Page %d]", currentPage),
			Items:     epList,
			Size:      10,
			Searcher:  searcher(epList),
			CursorPos: cursorPos,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}:",
				Active:   "\033[96m󰣇\033[0m {{ .Title | yellow }}",
				Inactive: "   {{ .Title | faint }}",
				Selected: "\033[96m󰣇\033[0m {{ .Title | green }}",
			},
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return
		}
		cursorPos = idx
		selected := epList[idx]
		switch selected.Title {
		case next:
			currentPage--
			continue
		case prev:
			currentPage++
			continue

		}
		url, err := selectResolution(selected.ID)
		if err != nil {
			fmt.Printf("✘ Error: %v", err)
			time.Sleep(5 * time.Second)
			return
		}
		if url == "Kembali" {
			continue
		}
		player.PlayVideo(url, selected.Title)
		utils.SaveState(aniID, selected.ID, selected.Title)
	}
}

// vim: ft=go
