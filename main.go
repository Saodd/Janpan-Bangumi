package main

import (
	"janpan-bangumi/libs"
)

func main() {
	libs.CheckWorkDir()
	filepaths := libs.RecurListMds("./bangumi/comment")
	years := libs.ParseComments(filepaths)
	libs.Upload(years)
}
