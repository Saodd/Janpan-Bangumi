package libs

import (
	"bytes"
	"log"
	"os"
	"regexp"
)

func HandleCommentFiles(files []string) (comments []*BangumiComment) {
	var reg, _ = regexp.Compile("(?s)#(.*?)\n(.*?)\n```yaml\n(.*?)\n```")
	for _, fp := range files {
		f, err := os.ReadFile(fp)
		if err != nil {
			log.Fatalln(err)
		}
		for _, part := range bytes.Split(f, []byte("----\n")) {
			part = bytes.TrimSpace(part)
			if len(part) == 0 {
				continue
			}
			subMatches := reg.FindAllSubmatch(part, -1)
			if len(subMatches) != 1 {
				log.Println("匹配错误！", string(part))
				continue
			}
			for _, sm := range subMatches {
				comment := &BangumiComment{
					Title:   string(bytes.TrimSpace(sm[1])),
					Content: string(bytes.TrimSpace(sm[2])),
					yaml:    sm[3],
				}
				comments = append(comments, comment)
			}
		}
	}
	return comments
}

type BangumiComment struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	yaml    []byte
}
