package torrentpotato

import (
	"encoding/json"
	"net/http"

	"github.com/cardigann/cardigann/torznab"
)

type Result struct {
	ReleaseName string `json:"release_name"`
	TorrentID   string `json:"torrent_id"`
	DetailsURL  string `json:"details_url"`
	DownloadURL string `json:"download_url"`
	ImdbID      string `json:"imdb_id"`
	Freeleech   bool   `json:"freeleech"`
	Type        string `json:"type"`
	Size        int    `json:"size"`
	Leechers    int    `json:"leechers"`
	Seeders     int    `json:"seeders"`
}

type Results struct {
	Results      []Result `json:"results"`
	TotalResults int      `json:"total_results"`
}

func Error(w http.ResponseWriter, err error) {
	b, err := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "  ")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadGateway)
	w.Write(append(b, '\n'))
}

func Output(w http.ResponseWriter, items []torznab.ResultItem) {
	results := []Result{}

	for _, item := range items {
		results = append(results, Result{
			ReleaseName: item.Title,
			TorrentID:   item.GUID,
			DetailsURL:  item.Comments,
			DownloadURL: item.Link,
			Type:        "movie",
			Size:        int(item.Size / 1024 / 1024),
			Leechers:    item.Peers - item.Seeders,
			Seeders:     item.Seeders,
		})
	}

	res := Results{results, len(results)}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		Error(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(append(b, '\n'))
}
