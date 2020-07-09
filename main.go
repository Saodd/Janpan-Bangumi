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
	if err := libs.SaveBlogs(blogs); err != nil {
		log.Fatalln(err)
	}
}
