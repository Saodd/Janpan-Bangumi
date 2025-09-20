package libs

import (
	"log"
	"os"

	"resty.dev/v3"
)

var (
	JulietBangumiUrl       string
	JulietBangumiPostToken string
)

func init() {
	if mode := os.Getenv("RUN_MODE"); mode == "" {
		JulietBangumiUrl = "http://localhost:20001"
	} else {
		JulietBangumiUrl = "https://api.lewinblog.com"
	}
	JulietBangumiPostToken = os.Getenv("JULIET_POST_TOKEN")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Upload(comments []*BangumiComment) {
	toUpload, toDelete := filterComments(comments)
	uploadComments(toUpload)
	deleteComments(toDelete)
}

func filterComments(comments []*BangumiComment) (toUpload, toDelete []*BangumiComment) {
	var serverHashMap = make(map[string]string) // 服务器端已有评论的 Title -> Hash 映射
	{
		var respBody struct {
			Code int               `json:"code"`
			Data []*BangumiComment `json:"data"` // 只有Title和Hash字段有值
		}
		client := resty.New()
		defer client.Close()
		resp, err := client.R().
			SetHeader("X-STAFF-TOKEN", JulietBangumiPostToken).
			SetResult(&respBody).
			Get(JulietBangumiUrl + "/bangumi/staff/list_hash")
		if err != nil {
			log.Fatalln(err)
		}
		if resp.IsError() {
			log.Fatalln("请求失败，状态码:", resp.StatusCode())
		}
		if respBody.Code != 0 {
			log.Fatalln("请求失败，Code=", respBody.Code)
		}
		for _, c := range respBody.Data {
			serverHashMap[c.Title] = c.Hash
		}
		log.Printf("已有评论哈希 %d 条\n", len(serverHashMap))
	}
	for _, c := range comments {
		if serverHashMap[c.Title] != c.Hash {
			toUpload = append(toUpload, c)
		}
	}
	// 找出服务器端有但本地没有的评论，准备删除
	var localHashMap = make(map[string]bool)
	for _, c := range comments {
		localHashMap[c.Title] = true
	}
	for title := range serverHashMap {
		if !localHashMap[title] {
			toDelete = append(toDelete, &BangumiComment{Title: title})
		}
	}
	return toUpload, toDelete
}

func uploadComments(comments []*BangumiComment) {
	const batchSize = 50 // 每次最多执行 50 条
	client := resty.New()
	defer client.Close()
	for i := 0; i < len(comments); i += batchSize {
		end := i + batchSize
		if end > len(comments) {
			end = len(comments)
		}
		var respBody struct {
			Code int `json:"code"`
		}
		resp, err := client.R().
			SetHeader("X-STAFF-TOKEN", JulietBangumiPostToken).
			SetBody(map[string]any{"comments": comments[i:end]}).
			SetResult(&respBody).
			Post(JulietBangumiUrl + "/bangumi/staff/upsert_comments")
		if err != nil {
			log.Fatalln(err)
		}
		if resp.IsError() {
			log.Fatalln("请求失败，状态码:", resp.StatusCode())
		}
		if respBody.Code != 0 {
			log.Fatalln("上传失败，Code=", respBody.Code)
		}
		log.Printf("上传 %d 条评论，进度 %d/%d\n", end-i, end, len(comments))
	}
}

func deleteComments(comments []*BangumiComment) {
	const batchSize = 50 // 每次最多执行 50 条
	client := resty.New()
	defer client.Close()
	for i := 0; i < len(comments); i += batchSize {
		end := i + batchSize
		if end > len(comments) {
			end = len(comments)
		}
		var respBody struct {
			Code int `json:"code"`
		}
		var titles []string
		for _, c := range comments[i:end] {
			titles = append(titles, c.Title)
		}
		resp, err := client.R().
			SetHeader("X-STAFF-TOKEN", JulietBangumiPostToken).
			SetBody(map[string]any{"titles": titles}).
			SetResult(&respBody).
			Post(JulietBangumiUrl + "/bangumi/staff/delete_comments")
		if err != nil {
			log.Fatalln(err)
		}
		if resp.IsError() {
			log.Fatalln("请求失败，状态码:", resp.StatusCode())
		}
		if respBody.Code != 0 {
			log.Fatalln("删除失败，Code=", respBody.Code)
		}
		log.Printf("删除 %d 条评论，进度 %d/%d\n", end-i, end, len(comments))
	}
}
