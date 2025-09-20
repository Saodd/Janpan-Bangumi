package libs

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

func ParseComments(filepaths []string) (comments []*BangumiComment) {
	// 1. 处理所有文件，拆分为多个 BangumiCommentRaw
	var raws []*BangumiCommentRaw
	var reg, _ = regexp.Compile("(?s)#(.*?)\n(.*?)\n```yaml\n(.*?)\n```")
	for _, fp := range filepaths {
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
				log.Fatalln("匹配错误！", string(part))
			}
			for _, sm := range subMatches {
				raw := &BangumiCommentRaw{
					Title:   string(bytes.TrimSpace(sm[1])),
					Content: string(bytes.TrimSpace(sm[2])),
					yaml:    sm[3],
				}
				raws = append(raws, raw)
			}
		}
	}
	// 2. 检查Title重复
	var titleMap = make(map[string]bool)
	for _, raw := range raws {
		if titleMap[raw.Title] {
			log.Fatalln("Title重复了:", raw.Title)
		}
		titleMap[raw.Title] = true
	}
	// 3. 转换为 BangumiComment
	for _, raw := range raws {
		var comment = &BangumiComment{}
		if err := yaml.Unmarshal(raw.yaml, comment); err != nil {
			log.Fatalln("yaml解析失败", "|", err, "|", string(raw.yaml))
		}
		comment.Title = raw.Title
		comment.Content = raw.Content
		comment.Hash = raw.calcHash()
		comment.check()
		comments = append(comments, comment)
	}
	return comments
}

func (c *BangumiCommentRaw) calcHash() string {
	hasher := md5.New()
	hasher.Write([]byte(c.Title))
	hasher.Write([]byte(c.Content))
	hasher.Write(c.yaml)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (c *BangumiComment) check() {
	// 1. 处理封面图
	if c.Cover == "" {
		fmt.Println("缺少 Cover:", c.Title)
	} else if strings.Contains(c.Cover, "hdslb.com") {
		// B站
		c.CoverS = c.Cover + "@200w_268h.webp"
	} else if strings.Contains(c.Cover, "/l/public/") {
		// 豆瓣大图
		c.CoverS = strings.ReplaceAll(c.Cover, "/l/public/", "/s_ratio_poster/public/")
	} else if strings.Contains(c.Cover, "/s_ratio_poster/public/") {
		// 豆瓣小图，处理粗心的时候
		c.CoverS = c.Cover
		c.Cover = strings.ReplaceAll(c.Cover, "/s_ratio_poster/public/", "/l/public/")
	} else if strings.Contains(c.Cover, "/s/public") {
		// 豆瓣图书
		c.CoverS = c.Cover
	} else if strings.Contains(c.Cover, "https://bookcover.yuewen.com/qdbimg/") {
		// 起点读书
		c.CoverS = c.Cover
	} else {
		fmt.Println("无法识别的 Cover:", c.Title)
		c.CoverS = c.Cover
	}
	// 2. 检查 YearMonth
	if c.YearMonth == 0 {
		fmt.Println("缺少 YearMonth:", c.Title)
	}
}

type BangumiCommentRaw struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	yaml    []byte
}

type BangumiComment struct {
	Title     string   `json:"title,omitempty" yaml:"title"`         // 番剧/作品标题
	Alias     []string `json:"alias,omitempty" yaml:"alias"`         // 别名
	YearMonth int      `json:"yearMonth,omitempty" yaml:"yearMonth"` // 首播时间，格式：202309
	Episode   string   `json:"episode,omitempty" yaml:"episode"`     // 集数
	Type      string   `json:"type,omitempty" yaml:"type"`           // 类型，bangumi/movie/other

	Link   string `json:"link,omitempty" yaml:"link"`     // 视频播放地址（B站）
	Cover  string `json:"cover,omitempty" yaml:"cover"`   // 封面URL
	CoverS string `json:"coverS" bson:"coverS"`           // 封面缩略图URL
	Douban string `json:"douban,omitempty" yaml:"douban"` // 豆瓣详情页URL

	MarkScore   int    `json:"markScore,omitempty" yaml:"markScore"`     // 评分，1-5
	MarkBrev    string `json:"markBrev,omitempty" yaml:"markBrev"`       // 简评
	MarkDate    string `json:"markDate,omitempty" yaml:"markDate"`       // 评分日期，格式：2023-09-01
	MarkEpisode string `json:"markEpisode,omitempty" yaml:"markEpisode"` // 评分时已经观看的集数
	Content     string `json:"content,omitempty"`                        // 详细评论，Markdown格式

	Hash string `json:"hash,omitempty"` // 属于哪一个20xx.md文件
}
