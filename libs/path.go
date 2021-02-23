package libs

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const ProjectDirname = "Janpan-Bangumi"

// CheckWorkDir 检查当前工作目录是否是项目根目录，不是的话就退出
func CheckWorkDir() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	if filepath.Base(workDir) != ProjectDirname {
		log.Fatalf("当前路径(%s)不是项目根目录(%s)！\n", filepath.Base(workDir), ProjectDirname)
	}
}

// RecurListMds 将递归遍历指定目录，返回所有 .md 文件的路径。
func RecurListMds(folder string) (mds []string) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			subFolder := filepath.Join(folder, file.Name())
			mds = append(mds, RecurListMds(subFolder)...)
		} else {
			if name := file.Name(); len(name) > 3 && name[len(name)-3:] == ".md" {
				mds = append(mds, filepath.Join(folder, file.Name()))
			}
		}
	}
	return
}
