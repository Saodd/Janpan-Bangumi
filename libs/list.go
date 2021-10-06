package libs

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func HandleListFiles(files []string) {
	var ch = make(chan []*BangumiData, 100)

	go func() {
		for _, fp := range files {
			f, err := os.ReadFile(fp)
			if err != nil {
				log.Fatalln(err)
			}
			var data []*BangumiData
			if err := json.Unmarshal(f, &data); err != nil {
				log.Fatalln(err)
			}
			ch <- data
		}
		close(ch)
	}()

	{
		req, _ := http.NewRequest("POST", JulietBangumiUrl+"/bangumi/drop-list", nil)
		req.Header.Set("X-STAFF-TOKEN", JulietBangumiPostToken)
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
	}
	for data := range ch {
		j, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", JulietBangumiUrl+"/bangumi/listV2", bytes.NewReader(j))
		req.Header.Set("X-STAFF-TOKEN", JulietBangumiPostToken)
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

type BangumiData struct {
	Cover       string   `json:"cover,omitempty"`
	Title       string   `json:"title,omitempty"`
	Link        string   `json:"link,omitempty"`
	YearMonth   int      `json:"yearMonth,omitempty"`
	Episode     string   `json:"episode,omitempty"`
	MarkStatus  int      `json:"markStatus,omitempty"`
	MarkScore   int      `json:"markScore,omitempty"`
	MarkBrev    string   `json:"markBrev,omitempty"`
	MarkDate    string   `json:"markDate,omitempty"`
	MarkEpisode string   `json:"markEpisode,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}
