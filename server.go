package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func GetUrl(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return "", err
	}

	req.Header.Add(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:52.9) Gecko/20100101 Goanna/3.4 Firefox/52.9 PaleMoon/27.6.0",
	)

	response, err := (&http.Client{Timeout: time.Duration(10) * time.Second}).Do(req)
	if nil != err {
		return "", err
	}

	defer func() {
		io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return "", err
	}

	return string(body), nil
}

func handlerYoutube(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req:   %+v\n", r)
	fmt.Println(r.URL.RequestURI())
	url := func() string {
		if "/you" == r.URL.RequestURI() {
			return "https://www.youtube.com/results?search_query=%E9%87%8E%E7%90%83"
		} else {
			return "https://wwww.youtube.com" + r.URL.RequestURI()
		}
	}()
	fmt.Println(url)
	body, _ := GetUrl(url)
	fmt.Fprintf(w, body)
}

func handlerAmazon(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req:   %+v\n", r)
	fmt.Println(r.URL.RequestURI())
	url := func() string {
		if "/ama" == r.URL.RequestURI() {
			return "https://www.amazon.co.jp/s/ref=nb_sb_noss_1?__mk_ja_JP&url=search-alias%3Dprime-instant-video&field-keywords=another"
		} else {
			return "https://www.amazon.co.jp" + r.URL.RequestURI()
			// return "https://www.amazon.co.jp/s/ref=nb_sb_noss_1?__mk_ja_JP&url=search-alias%3Dprime-instant-video&field-keywords=another"
		}
	}()
	fmt.Println(url)
	body, _ := GetUrl(url)
	fmt.Fprintf(w, body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Req:   %+v\n", r)
	fmt.Println(r.URL.RequestURI())
	url := func() string {
		if "/" == r.URL.RequestURI() {
			return "https://www.happyon.jp/search?q=another"
		} else {
			return "https://www.happyon.jp" + r.URL.RequestURI()
		}
	}()
	fmt.Println(url)
	body, _ := GetUrl(url)
	fmt.Fprintf(w, body)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/ama", handlerAmazon)
	http.HandleFunc("/you", handlerYoutube)
	http.ListenAndServe("localhost:12345", nil)
	r, _ := GetUrl("http://www.yahoo.co.jp")
	fmt.Println(r)
}
