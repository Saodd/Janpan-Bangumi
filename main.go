package main

import (
	"janpan-bangumi/libs"
	"log"
)

func main() {
	libs.CheckWorkDir()
	files := libs.RecurListMds("./blog")
	blogs, err := libs.ParseBlogFiles(files)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("成功读取%d篇影评。\n", len(blogs))
	if err := libs.PostBlogs(blogs); err != nil {
		log.Fatalln(err)
	}
}
