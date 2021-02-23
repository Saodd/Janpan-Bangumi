package libs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"
)

const (
	PublicBlogDir = "./public"
	blogInputExpr = "(?s)\n## (.*?)\n(.*?)\n```yaml lw-blog-meta(.*?)```"
)

var (
	blogPattern            *regexp.Regexp
	JulietBangumiUrl       string
	JulietBangumiPostToken string
)

type Blog struct {
	// 番剧原始信息
	BangumiName   string `json:"bangumi_name" yaml:"BangumiName"`
	BangumiDate   string `json:"bangumi_date" yaml:"BangumiDate"`
	BangumiSource string `json:"bangumi_source" yaml:"BangumiSource"`
	BangumiImage  string `json:"bangumi_image" yaml:"BangumiImage"`
	SeriesName    string `json:"series_name" yaml:"SeriesName"`
	SeriesCode    string `json:"series_code" yaml:"SeriesCode"`
	// 本人观看数据
	WatchDate    string `json:"watch_date" yaml:"WatchDate"`
	WatchEpisode string `json:"watch_episode" yaml:"WatchEpisode"`
	// 评论
	MarkBrev  string `json:"mark_brev" yaml:"MarkBrev"`
	MarkScore int    `json:"mark_score" yaml:"MarkScore"`
	MarkBody  string `json:"mark_body"`
}

func NewBlog(title, meta, body []byte) (*Blog, error) {
	var blog = &Blog{}
	if err := yaml.Unmarshal(meta, blog); err != nil {
		return nil, err
	}
	blog.MarkBody = string(body)

	if blog.BangumiName == "" {
		blog.BangumiName = string(title)
	}
	if blog.MarkBrev == "" {
		blog.MarkBrev = blog.MarkBody[:50] + " ..."
	}
	return blog, nil
}

func ParseBlogFiles(filePaths []string) (blogs []*Blog, err error) {
	for _, p := range filePaths {
		fileBlogs, err := parseBlogFile(p)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("解析失败(%s): %s\n", p, err))
		}
		blogs = append(blogs, fileBlogs...)
	}
	return blogs, nil
}

func parseBlogFile(filePath string) ([]*Blog, error) {
	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var blogs []*Blog
	subMatches := blogPattern.FindAllSubmatch(text, -1)
	for _, sm := range subMatches {
		blog, err := NewBlog(sm[1], sm[3], sm[2])
		if err != nil {
			log.Println(err)
		} else {
			blogs = append(blogs, blog)
		}
	}

	return blogs, nil
}

func SaveBlogs(blogs []*Blog) error {
	if err := os.MkdirAll(PublicBlogDir, 0755); err != nil {
		return err
	}
	body, err := json.Marshal(blogs)
	if err != nil {
		return err
	}
	fileName := path.Join(PublicBlogDir, "blogs.json")
	return ioutil.WriteFile(fileName, body, 0755)
}

type PostSegment struct {
	PostTs int64   `json:"post_ts"`
	Action int     `json:"action"` // 约定：0是未完结，1是完结信号。
	Blogs  []*Blog `json:"blogs"`
	Token  string  `json:"token"`
}

func PostBlogs(blogs []*Blog) error {
	if len(blogs) == 0 {
		return nil
	}
	// 构建推送数据结构体。用时间戳来做原子性。
	var seg = PostSegment{PostTs: time.Now().Unix(), Token: JulietBangumiPostToken}
	// 分段发送
	for i := 0; i < len(blogs); i += 10 {
		var right = i + 10
		if right >= len(blogs) {
			right = len(blogs)
			seg.Action = 1
		}
		seg.Blogs = blogs[i:right]
		if err := postBlogsSegment(seg); err != nil {
			return err
		}
	}
	return nil
}

func postBlogsSegment(seg PostSegment) error {
	body, _ := json.Marshal(seg)
	r := bytes.NewReader(body)
	resp, err := http.Post(JulietBangumiUrl, "", r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Post 失败：%s", resp.Status))
	}
	return nil
}

func init() {
	if mode := os.Getenv("RUN_MODE"); mode == "" {
		JulietBangumiUrl = "http://localhost:7777/bangumi/list"
	} else {
		JulietBangumiUrl = "https://api.lewinblog.com/bangumi/list"
	}
	JulietBangumiPostToken = os.Getenv("JULIET_POST_TOKEN")

	blogPattern, _ = regexp.Compile(blogInputExpr)
}
