package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	run(2019, 1)
	run(2019, 4)
	run(2019, 7)
	run(2019, 10)
}

func run(year, month int) {
	const pageNum = 1
	const pageSize = 100

	var u = fmt.Sprintf(
		`https://api.bilibili.com/pgc/season/index/result?season_version=-1&spoken_language_type=1&area=2&is_finish=-1&copyright=-1&season_status=-1&season_month=%d&year=%%5B%d%%2C%d)&style_id=-1&order=3&st=1&sort=0&page=%d&season_type=1&pagesize=%d&type=1`,
		month, year, year+1, pageNum, pageSize,
	)

	resp, err := http.Get(u)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data BilibiliResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalln(err)
	}
	if len(data.Data.List) >= pageSize {
		log.Fatalln("数量大于", pageSize)
	}

	var res []*BangumiData
	for _, item := range data.Data.List {
		if strings.Contains(item.Title, "配版") {
			continue
		}
		if strings.Contains(item.Title, "地區") {
			continue
		}
		if strings.Contains(item.Title, "僅限") {
			continue
		}

		li := &BangumiData{
			Cover:      strings.Replace(item.Cover, "http://", "https://", 1),
			Title:      item.Title,
			Link:       strings.Replace(item.Link, "http://", "https://", 1),
			YearMonth:  year*100 + month,
			MarkStatus: 0,
		}
		res = append(res, li)
	}

	fp := filepath.Join("bangumi/list/bilibili", fmt.Sprintf("%d.json", year*100+month))
	fd, _ := json.Marshal(res)
	if err := ioutil.WriteFile(fp, fd, 0755); err != nil {
		log.Fatalln(err)
	}
}

type BilibiliResponse struct {
	Code int `json:"code"`
	Data struct {
		HasNext int `json:"has_next"`
		List    []*struct {
			Badge     string `json:"badge"`
			BadgeInfo struct {
				BgColor      string `json:"bg_color"`
				BgColorNight string `json:"bg_color_night"`
				Text         string `json:"text"`
			} `json:"badge_info"`
			BadgeType  int    `json:"badge_type"`
			Cover      string `json:"cover"`
			IndexShow  string `json:"index_show"`
			IsFinish   int    `json:"is_finish"`
			Link       string `json:"link"`
			MediaId    int    `json:"media_id"`
			Order      string `json:"order"`
			OrderType  string `json:"order_type"`
			SeasonId   int    `json:"season_id"`
			SeasonType int    `json:"season_type"`
			Title      string `json:"title"`
			TitleIcon  string `json:"title_icon"`
		} `json:"list"`
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
}
type BangumiData struct {
	Cover      string `json:"cover,omitempty" yaml:"cover"`
	Title      string `json:"title,omitempty" yaml:"title"`
	Link       string `json:"link,omitempty" yaml:"link"`
	YearMonth  int    `json:"yearMonth,omitempty" yaml:"yearMonth"`
	MarkStatus int    `json:"markStatus,omitempty" yaml:"markStatus"`

	//Episode     string   `json:"episode,omitempty" yaml:"episode"`
	//MarkScore   int      `json:"markScore,omitempty" yaml:"markScore"`
	//MarkBrev    string   `json:"markBrev,omitempty" yaml:"markBrev"`
	//MarkDate    string   `json:"markDate,omitempty" yaml:"markDate"`
	//MarkEpisode string   `json:"markEpisode,omitempty" yaml:"markEpisode"`
	//Tags        []string `json:"tags,omitempty" yaml:"tags"`
}
