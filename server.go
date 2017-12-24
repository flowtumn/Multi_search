package muls

/*
* Step1: UserAccess  (e.g http://localhost:10101)
*  -> Response search only page.
* Step2: User submit.
*  -> Return the HTML with keyword embedded in URL. (e.g  http://localhost:10101/v/hulu?keyword=xxxxx)
*     If search type is "shop" to return raw url. (e.g http://www.1999.co.jp/search?serchkey=xxxxxx)
* Finish
 */

//検索Page -> Submit -> iframe -> {auto action}

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type SearchType int64

const (
	//対応するパラメーター
	KEYWORD = "keyword"
	SEARCH  = "search"

	//ParseするBaseとなるHtml
	INDEX_HTML  = "index.html"
	SEARCH_HTML = "search.html"
)

const (
	Video SearchType = 0
	Shop  SearchType = 1
)

type SearchProxyServer struct {
	_HomeDir    string
	_Router     []Router
	_Mux        *http.ServeMux
	_Server     *http.Server
	_Port       int
	_SupportUri map[string]Router
}

func (self *SearchProxyServer) _GetHttpServerUrl() string {
	return "http://" + self._Server.Addr
}

func (self *SearchProxyServer) _AtRouterByName(name string) *Router {
	//Nameを基に検索。
	for _, v := range self._Router {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func (self *SearchProxyServer) _AtRouter(r *http.Request) *Router {
	//登録されているrouter情報から取得。
	for _, v := range self._Router {
		if v.Endpoint == r.URL.Path {
			return &v
		}
	}
	return nil
}

func (self *SearchProxyServer) _AtSupportUrl(r *http.Request) *Router {
	//Referから判定
	for k, v := range self._SupportUri {
		refer := r.Referer()
		keyLength := len(k)

		if keyLength < len(refer) {
			refer = refer[:keyLength]
		}

		if k == refer {
			return &v
		}
	}
	return nil
}

/**
 * Index.htmlに対応するHandler
 */
func (self *SearchProxyServer) _Handle_Index(w http.ResponseWriter) {
	t, err := template.ParseFiles(self._HomeDir + "/" + INDEX_HTML)
	if nil != err {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	type Info struct {
		Host string
	}

	err = t.Execute(
		w,
		Info{
			Host: self._GetHttpServerUrl(),
		},
	)

	if nil != err {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

/**
 * 検索に対応するHandler
 */
func (self *SearchProxyServer) _Handle_Search(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(self._HomeDir + "/" + SEARCH_HTML)
	if nil != err {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	keyword := r.FormValue("keyword")
	searchType, err := strconv.ParseInt(r.FormValue("search"), 10, 64)
	if nil != err {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	type Info struct {
		Host     string
		Keyword  string
		Frame1   string
		Frame2   string
		Checked1 string
		Checked2 string
		Checked3 string
	}

	err = t.Execute(
		w,
		Info{
			Host:    self._GetHttpServerUrl(),
			Keyword: keyword,
			Frame1: func() string {
				if r := self._AtRouterByName(
					func() string {
						switch (SearchType)(searchType) {
						case Video:
							return "hulu"
						case Shop:
							return "yodobashi"
						}
						return ""
					}(),
				); r != nil {
					//Videoなら、本Serverを経由
					if Video == (SearchType)(searchType) {
						return self._GetHttpServerUrl() + r.CreateEndpointWithKeyword(keyword)
					} else {
						//それ以外なら転送先のURIを返す。
						return r.CreateSearchUri(keyword)
					}
				}
				return ""
			}(),
			Frame2: func() string {
				if r := self._AtRouterByName(
					func() string {
						switch (SearchType)(searchType) {
						case Video:
							return "amazon"
						case Shop:
							return "1999"
						}
						return ""
					}(),
				); r != nil {
					//Videoなら、本Serverを経由
					if Video == (SearchType)(searchType) {
						return self._GetHttpServerUrl() + r.CreateEndpointWithKeyword(keyword)
					} else {
						//それ以外なら転送先のURIを返す。
						return r.CreateSearchUri(keyword)
					}
				}
				return ""
			}(),
			Checked1: func() string {
				if Video == (SearchType)(searchType) {
					return "checked"
				}
				return ""
			}(),
			Checked2: func() string {
				if Shop == (SearchType)(searchType) {
					return "checked"
				}
				return ""
			}(),
			Checked3: func() string {
				return ""
			}(),
		},
	)

	if nil != err {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (self *SearchProxyServer) _Handler(w http.ResponseWriter, r *http.Request) {
	url := func() string {
		//Referから引けるのなら、二回目以降のアクセスなのでBaseURL + RequestURIを返す
		if mp := self._AtSupportUrl(r); nil != mp {
			return mp.BaseUrl + r.RequestURI
		}

		//Referから引けなければ、初回アクセスとして検索キーを含めたUriを返す
		if mp := self._AtRouter(r); nil != mp {
			//keywordは引いたら消去する。(転送先にはノイズになるため)
			keyword := r.FormValue("keyword")
			if 0 != len(keyword) {
				r.Form.Del("keyword")

				//Methodも変換。
				r.Method = mp.Method.ToString()

				return mp.CreateSearchUri(keyword)
			}
		}

		return ""
	}()

	//遷移するべきURLが無い。
	if 0 == len(url) {
		//検索キーワードが含まれているか？
		keyword := r.FormValue("keyword")
		if len(keyword) == 0 {
			//検索キーワードが無いなら、IndexHandlerで対応。
			self._Handle_Index(w)
			return
		}

		//SearchHandlerで対応。
		self._Handle_Search(w, r)
		return
	}

	//Accept-Encodingは削除。(archiveが落ちてくるので)
	r.Header.Del("Accept-Encoding")

	err := DoRequest(
		url,
		r,
		func(response *http.Response) error {
			body, err := ioutil.ReadAll(response.Body)
			if nil != err {
				return err
			}

			w.WriteHeader(response.StatusCode)
			w.Write(body)
			return nil
		},
	)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * Serverを初期化し、Routerを結び付けます。
 */
func (self *SearchProxyServer) _Initialize(homeDir string, port int) error {
	self._Server = &http.Server{}
	self._SupportUri = map[string]Router{}
	self._HomeDir = homeDir

	for _, v := range []Router{
		//Hulu
		Router{
			Name:      "hulu",
			Method:    GET,
			Endpoint:  "/v/hulu",
			BaseUrl:   "https://www.happyon.jp",
			SearchUrl: "/search?q=",
		},
		//Amazon
		Router{
			Name:      "amazon",
			Method:    GET,
			Endpoint:  "/v/amazon",
			BaseUrl:   "https://www.amazon.co.jp",
			SearchUrl: "/s/ref=nb_sb_noss_1?__mk_ja_JP&url=search-alias%3Dprime-instant-video&field-keywords=",
		},
		//Youtube
		Router{
			Name:      "youtube",
			Method:    GET,
			Endpoint:  "/v/youtube",
			BaseUrl:   "https://www.youtube.com",
			SearchUrl: "/results?search_query=",
		},
		//Yodobashi.
		Router{
			Name:      "yodobashi",
			Method:    GET,
			Endpoint:  "/s/yodobashi",
			BaseUrl:   "http://www.yodobashi.com",
			SearchUrl: "/?word=",
		},
		//Hobby
		Router{
			Name:      "1999",
			Method:    GET,
			Endpoint:  "/s/1999",
			BaseUrl:   "http://www.1999.co.jp",
			SearchUrl: "/search?searchkey=",
		},
	} {
		//対応するrouterに追加。
		self._Router = append(self._Router, v)

		//Referから引けるようにするため、対応するURLも生成。
		self._SupportUri[fmt.Sprintf("http://localhost:%d%s", port, v.Endpoint)] = v
	}

	self._Mux = http.NewServeMux()
	self._Mux.HandleFunc("/", self._Handler)
	http.DefaultServeMux = self._Mux
	return nil
}

func (self *SearchProxyServer) Listen(host string) error {
	self._Server.Addr = fmt.Sprintf("%s:%d", host, self._Port)
	return self._Server.ListenAndServe()
}

func (v SearchProxyServer) GetRouters() []Router {
	return v._Router
}

/**
 * 指定したURLに対してRequestを発行します。
 */
func DoRequest(
	url string,
	request *http.Request,
	callback func(*http.Response) error,
) error {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	b := ioutil.NopCloser(bytes.NewReader(body))
	req, err := http.NewRequest(request.Method, url, b)
	if nil != err {
		return err
	}

	//Copy Request-Header
	for k, v := range request.Header {
		req.Header[k] = v
	}

	response, err := (&http.Client{Timeout: time.Duration(10) * time.Second}).Do(req)
	if nil != err {
		return err
	}

	defer func() {
		io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	return callback(response)
}

/**
 * Serverを生成します。
 */
func CreateSearchProxyServer(homeDir string, port int) (*SearchProxyServer, error) {
	r := &SearchProxyServer{}
	r._Port = port
	if err := r._Initialize(homeDir, port); nil != err {
		return nil, err
	}
	return r, nil
}
