package libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	JulietBangumiUrl       string
	JulietBangumiPostToken string
)

func init() {
	if mode := os.Getenv("RUN_MODE"); mode == "" {
		JulietBangumiUrl = "http://localhost:7777"
	} else {
		JulietBangumiUrl = "https://api.lewinblog.com"
	}
	JulietBangumiPostToken = os.Getenv("JULIET_POST_TOKEN")
}

func Upload(datas []*BangumiData, comments []*BangumiComment) {
	{
		req, _ := http.NewRequest("POST", JulietBangumiUrl+"/bangumi/drop", nil)
		req.Header.Set("X-STAFF-TOKEN", JulietBangumiPostToken)
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
	}
	for i := 0; i < len(datas); i += 500 {
		part := datas[i:]
		if len(part) > 500 {
			part = part[:500]
		}
		j, _ := json.Marshal(part)
		req, _ := http.NewRequest("POST", JulietBangumiUrl+"/bangumi/listV2", bytes.NewReader(j))
		req.Header.Set("X-STAFF-TOKEN", JulietBangumiPostToken)
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
	}
	for i := 0; i < len(comments); i += 100 {
		part := comments[i:]
		if len(part) > 100 {
			part = part[:100]
		}
		j, _ := json.Marshal(part)
		req, _ := http.NewRequest("POST", JulietBangumiUrl+"/bangumi/comment", bytes.NewReader(j))
		req.Header.Set("X-STAFF-TOKEN", JulietBangumiPostToken)
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func CombineListAndComment(datas []*BangumiData, comments []*BangumiComment) ([]*BangumiData, []*BangumiComment) {
	var mp = make(map[string]*BangumiData)
	for _, data := range datas {
		if mp[data.Title] != nil {
			log.Fatalln("List重复了:", data.Title)
		}
		mp[data.Title] = data
	}
	var comments2 []*BangumiComment
	var datas2 []*BangumiData
	for _, comment := range comments {
		data := mp[comment.Title]
		if data == nil {
			//log.Println("没有对应的List:", comment.Title)
			data = &BangumiData{Title: comment.Title}
			mp[comment.Title] = data
		}
		if err := yaml.Unmarshal(comment.yaml, data); err != nil {
			log.Fatalln("yaml解析失败", "|", err, "|", string(comment.yaml))
		}
		comments2 = append(comments2, comment)
		datas2 = append(datas2, data)
	}

	return datas2, comments2
}

func CheckDatas(datas []*BangumiData) {
	for _, data := range datas {
		if data.Cover == "" {
			fmt.Println("缺少 Cover:", data.Title)
		} else if strings.Contains(data.Cover, "hdslb.com") {
			// B站
			data.CoverS = data.Cover + "@200w_268h.webp"
		} else if strings.Contains(data.Cover, "/l/public/") {
			// 豆瓣大图
			data.CoverS = strings.ReplaceAll(data.Cover, "/l/public/", "/s_ratio_poster/public/")
		} else if strings.Contains(data.Cover, "/s_ratio_poster/public/") {
			// 豆瓣小图，处理粗心的时候
			data.CoverS = data.Cover
			data.Cover = strings.ReplaceAll(data.Cover, "/s_ratio_poster/public/", "/l/public/")
		} else if strings.Contains(data.Cover, "/s/public") {
			// 豆瓣图书
			data.CoverS = data.Cover
		} else if strings.Contains(data.Cover, "https://bookcover.yuewen.com/qdbimg/") {
			// 起点读书
			data.CoverS = data.Cover
		} else {
			fmt.Println("无法识别的 Cover:", data.Title)
			data.CoverS = data.Cover
		}

		if data.YearMonth == 0 {
			fmt.Println("缺少 YearMonth:", data.Title)
		}
	}
}
