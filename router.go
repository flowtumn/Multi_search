package muls

import "net/url"

const (
	Router_Keyword = "keyword"
)

type Router struct {
	Name      string
	Endpoint  string
	BaseUrl   string
	SearchUrl string
}

/**
 * 検索ワードを含めたURIを生成。
 */
func (v Router) CreateSearchUri(keyword string) string {
	return v.BaseUrl + v.SearchUrl + url.QueryEscape(keyword)
}

/**
 * keywordを含めたEndpointを生成。
 */
func (v Router) CreateEndpointWithKeyword(keyword string) string {
	return v.Endpoint + "?keyword=" + keyword
}
