package main

import (
	"janpan-bangumi/libs"
)

func main() {
	libs.CheckWorkDir()
	listFiles := libs.RecurListJsons("./bangumi/list")
	libs.HandleListFiles(listFiles)

	//blogs, err := libs.ParseBlogFiles(listFiles)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Printf("成功读取%d篇影评。\n", len(blogs))
	//libs.QuickSortBlog(blogs)
	//if err := libs.PostBlogs(blogs); err != nil {
	//	log.Fatalln(err)
	//}
}
