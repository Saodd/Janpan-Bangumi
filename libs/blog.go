package libs

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
)

const PublicBlogDir = "./public"
const blogInputExpr = "(?s)## (.*?)\n(.*?)\n```yaml lw-blog-meta(.*?)```"

var blogPattern *regexp.Regexp

func init() {
	blogPattern, _ = regexp.Compile(blogInputExpr)
}

type Blog struct {
	// 番剧原始信息
	BangumiName   string `json:"bangumi_name" yaml:"BangumiName"`
	BangumiDate   string `json:"bangumi_date" yaml:"BangumiDate"`
	BangumiSource string `json:"bangumi_source" yaml:"BangumiSource"`
	BangumiImage  string `json:"bangumi_image" yaml:"BangumiImage"`
	SeriesName    string `json:"series_name" yaml:"SeriesName"`
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
	if blog.SeriesName == "" {
		blog.SeriesName = blog.BangumiName
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
