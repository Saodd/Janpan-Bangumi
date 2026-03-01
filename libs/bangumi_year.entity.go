package libs

type BangumiYearReq struct {
	Year     string           `json:"year"`
	Md5      string           `json:"md5"`
	Bangumis []*BangumiServer `json:"bangumis"`
}
