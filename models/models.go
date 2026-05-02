package models

type Episodes struct {
	Data struct {
		Episode []Episode `json:"episode"`
	} `json:"data"`
}

type Episode struct {
	ID      string `json:"id"`
	EpTitle string `json:"title"`
	Index   string `json:"index"`
	// Image        string `json:"image"`
}

type Detail struct {
	Data struct {
		MovieData Movie `json:"movie"`
	} `json:"data"`
}

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	/*
		Description string `json:"synopsis"`
		AltTitle    string `json:"synonyms"`
		Genre       string `json:"genre"`
	*/
}

// search
type Movies struct {
	Data struct {
		Movie []Movie `json:"movie"`
	} `json:"data"`
}

type ServerResponse struct {
	Data struct {
		Server []Server `json:"server"`
	} `json:"data"`
}

type Server struct {
	ID      string `json:"id"`
	Link    string `json:"link"`
	Quality string `json:"quality"`
	Type    string `json:"type"`
	// FileSize string `json:"key_file_size"`
}

type EpisodeResult struct {
	ID      string
	EpTitle string
}

// vim: ft=go
