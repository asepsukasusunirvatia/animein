package models

type SearchResponse struct {
	Data struct {
		Movie []Movie `json:"movie"`
	} `json:"data"`
}

type DetailResponse struct {
	Data struct {
		MovieData Movie `json:"movie"`
	} `json:"data"`
}

type EpisodesResponse struct {
	Data struct {
		Episode []Episode `json:"episode"`
	} `json:"data"`
}

type ServerResponse struct {
	Data struct {
		Server []Server `json:"server"`
	} `json:"data"`
}

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Genre string `json:"genre"`
	// Description string `json:"synopsis,omitempty"`
	// AltTitle    string `json:"synonyms,omitempty"`
}

type Episode struct {
	ID      string `json:"id"`
	EpTitle string `json:"title"`
	Index   string `json:"index"`
	// Image string `json:"image,omitempty"`
}

type Server struct {
	ID       string `json:"id"`
	Link     string `json:"link"`
	Quality  string `json:"quality"`
	Type     string `json:"type"`
	FileSize string `json:"key_file_size,omitempty"`
}

type EpisodeResult struct {
	ID      string
	EpTitle string
}

type Dict map[string]string

// vim: ft=go
