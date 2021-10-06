package libs

import (
	"encoding/json"
	"log"
	"os"
)

func HandleListFiles(files []string) (datas []*BangumiData) {
	for _, fp := range files {
		f, err := os.ReadFile(fp)
		if err != nil {
			log.Fatalln(err)
		}
		var data []*BangumiData
		if err := json.Unmarshal(f, &data); err != nil {
			log.Fatalln(err)
		}
		datas = append(datas, data...)
	}
	return
}

type BangumiData struct {
	Cover      string `json:"cover,omitempty"`
	Title      string `json:"title,omitempty"`
	Link       string `json:"link,omitempty"`
	YearMonth  int    `json:"yearMonth,omitempty"`
	MarkStatus int    `json:"markStatus,omitempty"`

	Episode     string   `json:"episode,omitempty" yaml:"episode"`
	MarkScore   int      `json:"markScore,omitempty" yaml:"markScore"`
	MarkBrev    string   `json:"markBrev,omitempty" yaml:"markBrev"`
	MarkDate    string   `json:"markDate,omitempty" yaml:"markDate"`
	MarkEpisode string   `json:"markEpisode,omitempty" yaml:"markEpisode"`
	Tags        []string `json:"tags,omitempty" yaml:"tags"`
}
