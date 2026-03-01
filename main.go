package main

import (
	"janpan-bangumi/libs"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 配置
	endpoint := "http://localhost:20001"
	if val := os.Getenv("JULIET_ENDPOINT"); val != "" {
		endpoint = val
	}
	staffToken := os.Getenv("JULIET_POST_TOKEN")

	// 文件解析
	fm := libs.NewFileManager()
	yearReqs, err := fm.ReadAndParse()
	if err != nil {
		log.Fatalln("解析失败:", err)
	}

	// 上传
	uploader := libs.NewUploader(endpoint, staffToken)
	if err := uploader.Upload(yearReqs); err != nil {
		log.Fatalln("上传失败:", err)
	}
}
