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
	Title     string   `json:"title,omitempty" yaml:"title"`
	Alias     []string `json:"alias,omitempty" yaml:"alias"`
	YearMonth int      `json:"yearMonth,omitempty" yaml:"yearMonth"`
	Episode   string   `json:"episode,omitempty" yaml:"episode"`

	Link   string `json:"link,omitempty" yaml:"link"`     // 视频播放地址（B站）
	Cover  string `json:"cover,omitempty" yaml:"cover"`   // 封面URL
	Douban string `json:"douban,omitempty" yaml:"douban"` // 豆瓣详情页URL

	MarkScore   int    `json:"markScore,omitempty" yaml:"markScore"`
	MarkBrev    string `json:"markBrev,omitempty" yaml:"markBrev"`
	MarkDate    string `json:"markDate,omitempty" yaml:"markDate"`
	MarkEpisode string `json:"markEpisode,omitempty" yaml:"markEpisode"`

	Tags []string `json:"tags,omitempty" yaml:"tags"`
	Type string   `json:"type,omitempty" yaml:"type"`
}
