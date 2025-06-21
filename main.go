package main

import (
	"janpan-bangumi/libs"
	"log"
)

func main() {
	libs.CheckWorkDir()
	commentFiles := libs.RecurListMds("./bangumi/comment")
	comments := libs.HandleCommentFiles(commentFiles)

	datas, comments := libs.CombineListAndComment(nil, comments)
	libs.CheckDatas(datas)
	log.Printf("共计：%d个Comment", len(comments))
	libs.Upload(datas, comments)
}
