package libs

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

type BangumiLocal struct {
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
	Content     string `json:"content,omitempty" yaml:"-"`               // 详细评论，Markdown格式
}

func (b *BangumiLocal) Parse(title, content, yamlData string) error {
	// 0. 解析 YAML 数据
	if err := yaml.Unmarshal([]byte(yamlData), &b); err != nil {
		return fmt.Errorf("yaml解析失败 | %v | %s\n", err, yamlData)
	}
	b.Title = title
	b.Content = content

	// 1. 处理封面图
	if b.Cover == "" {
		log.Println("缺少 Cover:", b.Title)
	} else if strings.Contains(b.Cover, "hdslb.com") {
		// B站
		b.CoverS = b.Cover + "@200w_268h.webp"
	} else if strings.Contains(b.Cover, "doubanio.com") {
		if b.Douban == "" {
			// 豆瓣图片需要Referer头，后续下载时会用到
			log.Println("缺少 Douban:", b.Title)
		}
		if strings.Contains(b.Cover, "/l/public/") {
			// 豆瓣大图
			b.CoverS = strings.ReplaceAll(b.Cover, "/l/public/", "/s_ratio_poster/public/")
		} else if strings.Contains(b.Cover, "/s_ratio_poster/public/") {
			// 豆瓣小图，处理粗心的时候
			b.CoverS = b.Cover
			b.Cover = strings.ReplaceAll(b.Cover, "/s_ratio_poster/public/", "/l/public/")
		} else if strings.Contains(b.Cover, "/s/public") {
			// 豆瓣图书
			b.CoverS = b.Cover
		}
	} else if strings.Contains(b.Cover, "https://bookcover.yuewen.com/qdbimg/") {
		// 起点读书
		b.CoverS = b.Cover
	} else {
		b.CoverS = b.Cover
	}

	// 2. 检查 YearMonth
	if b.YearMonth == 0 {
		return fmt.Errorf("缺少 YearMonth: %s", b.Title)
	}
	return nil
}

func (b *BangumiLocal) AsServer() *BangumiServer {
	return &BangumiServer{
		Title:     b.Title,
		Alias:     b.Alias,
		YearMonth: b.YearMonth,
		Episode:   b.Episode,
		Type:      b.Type,

		Link:     b.Link,
		Cover:    b.CoverS, // 一律使用缩略图
		CoverMd5: fmt.Sprintf("%x", md5.Sum([]byte(b.Cover))),
		Douban:   b.Douban,

		MarkScore:   b.MarkScore,
		MarkBrev:    b.MarkBrev,
		MarkDate:    b.MarkDate,
		MarkEpisode: b.MarkEpisode,
		Content:     b.Content,
	}
}

type BangumiServer struct {
	Title     string   `json:"title" bson:"_id"`           // 番剧/作品标题
	Alias     []string `json:"alias" bson:"alias"`         // 别名
	YearMonth int      `json:"yearMonth" bson:"yearMonth"` // 首播时间，格式：202309
	Episode   string   `json:"episode" bson:"episode"`     // 总集数
	Type      string   `json:"type" bson:"type"`           // 类型，bangumi/movie/other

	Link     string `json:"link" bson:"link"`         // 视频播放地址（B站）
	Cover    string `json:"cover" bson:"cover"`       // 封面URL（统一使用缩略图）
	CoverMd5 string `json:"coverMd5" bson:"coverMd5"` // 封面URL的MD5值（可用于从/oss读取）
	Douban   string `json:"douban" bson:"douban"`     // 豆瓣详情页URL

	MarkScore   int    `json:"markScore" bson:"markScore"`     // 评分，1-5
	MarkBrev    string `json:"markBrev" bson:"markBrev"`       // 简评
	MarkDate    string `json:"markDate" bson:"markDate"`       // 评分日期，格式：2023-09-01
	MarkEpisode string `json:"markEpisode" bson:"markEpisode"` // 已经观看的集数
	Content     string `json:"content" bson:"content"`         // 详细评论，Markdown格式
}

func (b *BangumiServer) AsLocal() *BangumiLocal {
	return &BangumiLocal{
		Title:     b.Title,
		Alias:     b.Alias,
		YearMonth: b.YearMonth,
		Episode:   b.Episode,
		Type:      b.Type,

		Link:   b.Link,
		Cover:  b.Cover,
		CoverS: b.Cover,
		Douban: b.Douban,

		MarkScore:   b.MarkScore,
		MarkBrev:    b.MarkBrev,
		MarkDate:    b.MarkDate,
		MarkEpisode: b.MarkEpisode,
		Content:     b.Content,
	}
}
