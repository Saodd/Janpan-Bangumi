package libs

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileManager struct {
	dirPath string // 目录路径
}

func NewFileManager() *FileManager {
	return &FileManager{dirPath: "./bangumi"}
}

// RecurListMds 将递归遍历指定目录，返回所有 .md 文件的路径。
func (fm *FileManager) recurListMds(folder string) (mds []string) {
	files, _ := os.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			subFolder := filepath.Join(folder, file.Name())
			mds = append(mds, fm.recurListMds(subFolder)...)
		} else {
			if strings.HasSuffix(file.Name(), ".md") {
				mds = append(mds, filepath.Join(folder, file.Name()))
			}
		}
	}
	return
}

func (fm *FileManager) ReadAndParse() ([]BangumiYearReq, error) {
	// 递归遍历目录下的所有 .md 文件，读取内容并解析成结构化数据
	var reqs []BangumiYearReq
	filepaths := fm.recurListMds(fm.dirPath)
	var reg, _ = regexp.Compile("(?s)#(.*?)\n(.*?)\n```yaml\n(.*?)\n```")

	for _, fp := range filepaths {
		f, err := os.ReadFile(fp)
		if err != nil {
			return nil, err
		}

		yearName := strings.TrimSuffix(filepath.Base(fp), ".md")
		yearMd5 := fmt.Sprintf("%x", md5.Sum(f)) // 计算整个年份文件的 MD5

		var servers []*BangumiServer

		for _, part := range bytes.Split(f, []byte("----\n")) {
			part = bytes.TrimSpace(part)
			if len(part) == 0 {
				continue
			}
			subMatches := reg.FindAllSubmatch(part, -1)
			if len(subMatches) != 1 {
				return nil, fmt.Errorf("匹配错误！ %s", string(part))
			}
			for _, sm := range subMatches {
				title := string(bytes.TrimSpace(sm[1]))
				content := string(bytes.TrimSpace(sm[2]))
				yamlData := string(sm[3])

				var local BangumiLocal
				if err := local.Parse(title, content, yamlData); err != nil {
					return nil, err
				}
				servers = append(servers, local.AsServer())
			}
		}

		reqs = append(reqs, BangumiYearReq{
			Year:     yearName,
			Md5:      yearMd5,
			Bangumis: servers,
		})
	}
	return reqs, nil
}

func (fm *FileManager) CheckDuplicated(reqs []BangumiYearReq) {
	titleMap := make(map[string]*BangumiServer)
	coverMap := make(map[string]*BangumiServer)
	coverMd5Map := make(map[string]string) // md5 -> cover
	for _, req := range reqs {
		for _, b := range req.Bangumis {
			if exist, ok := titleMap[b.Title]; ok {
				log.Printf("Title重复： %s | %s\n", b.Title, exist.Title)
			} else {
				titleMap[b.Title] = b
			}
			if exist, ok := coverMap[b.Cover]; ok {
				log.Printf("Cover重复： %s | %s\n", b.Title, exist.Title)
			} else {
				coverMap[b.Cover] = b
			}
			if existCover, ok := coverMd5Map[b.CoverMd5]; ok && existCover != b.Cover {
				log.Printf("罕见-URL的MD5重复了： %s | %s\n", b.Cover, existCover)
			}
		}
	}
}

func (fm *FileManager) MarshalBangumiLocal(b *BangumiLocal) (string, error) {
	yamlData, err := yaml.Marshal(b)
	if err != nil {
		return "", fmt.Errorf("yaml序列化失败: %v", err)
	}
	str := fmt.Sprintf(
		`
----
# %s

%s

%s
`, b.Title, b.Content, string(yamlData))
	return str, nil
}

func (fm *FileManager) SaveToFiles(yearReqs []BangumiYearReq) error {
	// 将结构化的数据按照格式写回文件（覆盖原文件）
	for _, req := range yearReqs {
		fp := filepath.Join(fm.dirPath, req.Year+".md")

		var parts []string
		for _, b := range req.Bangumis {
			local := b.AsLocal()
			// 对于 yaml 输出，需要先清除 title 和 content ，因为它们已经在正文中了
			local.Title = ""
			local.Content = ""
			// 转换字符串
			part, err := fm.MarshalBangumiLocal(local)
			if err != nil {
				return err
			}
			parts = append(parts, part)
		}

		if err := os.WriteFile(fp, []byte(strings.Join(parts, "")), 0644); err != nil {
			return err
		}
		log.Printf("成功保存文件: %s\n", fp)
	}
	return nil
}
