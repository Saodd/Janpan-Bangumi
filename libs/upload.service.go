package libs

import (
	"fmt"
	"log"

	"resty.dev/v3"
)

type Uploader struct {
	endpoint   string
	staffToken string
}

func NewUploader(endpoint, staffToken string) *Uploader {
	return &Uploader{
		endpoint:   endpoint,
		staffToken: staffToken,
	}
}

type BangumiYear struct {
	Year string `json:"year"`
	Md5  string `json:"md5"`
}

func (u *Uploader) Upload(yearReqs []BangumiYearReq) error {
	client := resty.New()
	defer client.Close()

	// 1. 查询服务器的现有数据，获取每个年份的MD5值
	var respBody struct {
		Code int           `json:"code"`
		Data []BangumiYear `json:"data"`
	}

	resp, err := client.R().
		SetHeader("X-STAFF-TOKEN", u.staffToken).
		SetResult(&respBody).
		Get(u.endpoint + "/bangumi/staff/years")
	if err != nil {
		return fmt.Errorf("获取服务器年份失败: %v", err)
	}
	if resp.IsError() {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode())
	}
	if respBody.Code != 0 {
		return fmt.Errorf("获取服务器年份失败，Code=%d", respBody.Code)
	}

	serverMd5Map := make(map[string]string)
	for _, y := range respBody.Data {
		serverMd5Map[y.Year] = y.Md5
	}

	log.Printf("服务器已有数据年份数: %d\n", len(serverMd5Map))

	// 2. 对比本地数据和服务器数据，找出需要上传的年份数据
	var toUpload []BangumiYearReq
	for _, req := range yearReqs {
		if serverMd5Map[req.Year] != req.Md5 {
			toUpload = append(toUpload, req)
		}
	}
	log.Printf("需要上传%d个年份文件\n", len(toUpload))

	// 3. 将需要上传的数据发送到服务器，更新服务器数据
	for _, req := range toUpload {
		var updateResp struct {
			Code int `json:"code"`
		}
		resp, err := client.R().
			SetHeader("X-STAFF-TOKEN", u.staffToken).
			SetBody(req).
			SetResult(&updateResp).
			Post(u.endpoint + "/bangumi/staff/update_year")

		if err != nil {
			return fmt.Errorf("上传年份数据失败 (%s): %v", req.Year, err)
		}
		if resp.IsError() {
			return fmt.Errorf("请求失败 (%s)，状态码: %d", req.Year, resp.StatusCode())
		}
		if updateResp.Code != 0 {
			return fmt.Errorf("上传失败 (%s)，Code=%d", req.Year, updateResp.Code)
		}
		log.Printf("成功上传年份数据: %s, 共 %d 部番剧\n", req.Year, len(req.Bangumis))
	}

	return nil
}
