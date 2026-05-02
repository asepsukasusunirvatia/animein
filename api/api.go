package api

import (
	"codeberg.org/Asep5K/animein/models"
	"codeberg.org/Asep5K/animein/utils"
	"encoding/json"
	"fmt"
	"sync"
)

var (
	// episodeCache,streamCache = make(map[string][]models.Episode), make(map[string][]models.Server)
	streamCache  = make(map[string][]models.Server)
	episodeCache = make(map[string][]models.Episode)
	cacheLock    sync.RWMutex
)

const (
	BaseURL      = "https://animeinweb.com/"
	baseEndpoint = BaseURL + "/api/proxy/3/2/"
)

func reqEpsPage(aniID string, pageNum int) ([]models.Episode, error) {
	targetURL := fmt.Sprintf("%smovie/episode/%s?page=%d", baseEndpoint, aniID, pageNum)
	res, err := utils.Request(targetURL)
	if err != nil {
		return nil, fmt.Errorf("GetEpisodesPage request failed: %w", err)
	}
	defer res.Body.Close()

	var eps models.Episodes
	if err := json.NewDecoder(res.Body).Decode(&eps); err != nil {
		return nil, fmt.Errorf("GetEpisodesPage decode failed: %w", err)
	}
	return eps.Data.Episode, nil
}

func GetEpisodesCached(aniID string, pageNum int) ([]models.Episode, error) {
	cacheKey := fmt.Sprintf("%s-%d", aniID, pageNum)
	cacheLock.RLock()
	data, found := episodeCache[cacheKey]
	cacheLock.RUnlock()

	if found {
		return data, nil
	}

	data, err := reqEpsPage(aniID, pageNum)
	if err != nil {
		return nil, err
	}
	cacheLock.Lock()
	episodeCache[cacheKey] = data
	cacheLock.Unlock()
	return data, nil
}

func GetPageCount(aniID string) (int, error) {
	eps, err := GetEpisodesCached(aniID, 0)
	if err != nil {
		return 0, fmt.Errorf("GetPageCount -> %w", err)
	}
	if len(eps) == 0 {
		return 0, fmt.Errorf("\033[9mid\033: '%s' belum rilis!\033[0m", aniID)
	}
	lastEp := utils.StrToInt(eps[0].Index)
	if lastEp <= 30 {
		return 0, nil
	}
	return (lastEp+30)/30 - 1, nil
}

func ParseEpisodes(episodes []models.Episode) <-chan models.EpisodeResult {
	ch := make(chan models.EpisodeResult)
	go func() {
		defer close(ch)
		for _, ep := range episodes {
			ch <- models.EpisodeResult{
				ID:      ep.ID,
				EpTitle: ep.EpTitle,
			}
		}
	}()
	return ch
}

func reqEpsInfo(epID string) []models.Server {
	targetURL := fmt.Sprintf("%sepisode/streamnew/%s", baseEndpoint, epID)
	res, err := utils.Request(targetURL)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	var info models.ServerResponse
	json.NewDecoder(res.Body).Decode(&info)
	return info.Data.Server
}

func GetEpsInfo(epId string) []models.Server {
	cacheLock.RLock()
	data, found := streamCache[epId]
	cacheLock.RUnlock()

	if found {
		return data
	}

	epsInfo := reqEpsInfo(epId)
	cacheLock.Lock()
	streamCache[epId] = epsInfo
	cacheLock.Unlock()
	return epsInfo
}

func SearchAnime(keyWord string) ([]models.Movie, error) {
	stop := utils.Loading("Mencari " + keyWord)
	res, err := utils.SearchRequest(keyWord)
	if err != nil {
		stop <- true
		return nil, err
	}
	defer res.Body.Close()

	var info models.Movies
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		stop <- true
		return nil, err
	}
	if len(info.Data.Movie) == 0 {
		stop <- true
		return nil, fmt.Errorf("Tidak menemukan: '%s'", keyWord)
	}
	stop <- true
	return info.Data.Movie, nil
}

// vim: ft=go
