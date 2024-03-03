package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	inputDir  = "C://Users/Lewin/Downloads"
	outputDir = "C://Users/Lewin/追番"
)

var (
	_buffer [512]byte
)

func main() {
	for true {
		run()
	}
}

func run() {
	dirs0 := RecurListDirs(inputDir)

	fmt.Println("----------开始运行----------")
	fmt.Println("请输入文件夹关键字：")
	var dirs []string
	var newDir string

	for {
		kw := ReadKeywords()
		if len(kw) == 0 {
			break
		} else {
			fmt.Printf("输入内容：[%s]，长度：%d\n", kw, len(kw))
		}

		dirs = nil
		fmt.Println("筛选结果：")
		for _, dir := range dirs0 {
			if strings.Index(filepath.Base(dir), kw) >= 0 {
				dirs = append(dirs, dir)
				fmt.Println(dir)
			}
		}
		newDir = filepath.Join(outputDir, kw)
		fmt.Printf("共 %d 条。\n", len(dirs))
		fmt.Printf("上述文件夹中的文件将被收集到[%s]，按回车键确认上述结果，或者重新输入文件夹关键字：\n", newDir)
	}

	err := os.Mkdir(newDir, 0755)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("已创建目录：", newDir)
	}
	for _, dir := range dirs {
		files, _ := ioutil.ReadDir(dir)
		for _, file := range files {
			if !file.IsDir() {
				name := file.Name()
				oldPath := filepath.Join(dir, name)
				newPath := filepath.Join(newDir, name)
				err := os.Rename(oldPath, newPath)
				if err != nil {
					log.Fatalln(err)
				} else {
					fmt.Printf("移动 【%s】 到 【%s】 \n", oldPath, newPath)
				}
			}
		}
	}
}

func RecurListDirs(folder string) (dirs []string) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			subFolder := filepath.Join(folder, file.Name())
			dirs = append(dirs, subFolder)
		}
	}
	return
}

func ReadKeywords() (keywords string) {
	n, err := os.Stdin.Read(_buffer[:])
	if err != nil {
		log.Fatalln(err)
	}
	keywords = string(_buffer[:n])
	keywords = strings.TrimSpace(keywords)
	return keywords
}
