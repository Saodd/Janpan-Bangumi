package libs

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
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
