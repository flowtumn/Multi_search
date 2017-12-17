package muls

import "net/url"

const (
	Router_Keyword = "keyword"
)

type HttpMethod int

const (
	GET HttpMethod = iota + 1
	POST
)

type Router struct {
	Name      string
	Method    HttpMethod
	Endpoint  string
	BaseUrl   string
	SearchUrl string
}

func (v HttpMethod) ToString() string {
	switch v {
	case GET:
		return "GET"
	case POST:
		return "POST"
	}
	panic("Unknown HttpMethod.")
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
