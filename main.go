package main

import (
	"janpan-bangumi/libs"
	"log"
)

func main() {
	libs.CheckWorkDir()
	listFiles := libs.RecurListJsons("./bangumi/list")
	datas := libs.HandleListFiles(listFiles)
	commentFiles := libs.RecurListMds("./bangumi/comment")
	comments := libs.HandleCommentFiles(commentFiles)

	datas, comments = libs.CombineListAndComment(datas, comments)
	log.Printf("共计：%d个List，%d个Comment", len(datas), len(comments))
	libs.Upload(datas, comments)
}
