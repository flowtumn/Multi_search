package muls

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_Router_CreateSearchUri(t *testing.T) {
	testValue := Router{
		Name:      "test",
		Endpoint:  "/a/b/c",
		BaseUrl:   "http://localhost:11111",
		SearchUrl: "/?s=",
	}

	for _, keyword := range []string{
		"Test",
		"てすと",
		"㈱ test　てす",
	} {
		expectedUri := fmt.Sprintf("%s%s%s", testValue.BaseUrl, testValue.SearchUrl, url.QueryEscape(keyword))
		if expectedUri != testValue.CreateSearchUri(keyword) {
			t.Fatalf("Doesn't match result url on CreateSearchUrl.")
		}
	}
}

func Test_Router_CreateEndpointWithKeyword(t *testing.T) {
	testValue := Router{
		Name:      "test",
		Endpoint:  "/a/b/c",
		BaseUrl:   "http://localhost:11111",
		SearchUrl: "/?s=",
	}

	for _, keyword := range []string{
		"Test",
		"てすと",
		"㈱ test　てす",
	} {
		expectedUri := fmt.Sprintf("%s?%s=%s", testValue.Endpoint, Router_Keyword, keyword)
		if expectedUri != testValue.CreateEndpointWithKeyword(keyword) {
			t.Fatalf("Doesn't match result url on CreateEndpointWithKeyword.")
		}
	}
}
