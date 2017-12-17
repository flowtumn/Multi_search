package muls

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Request内容をEchoするHandler
func _Echo_Handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//HeaderもEchoする
	for hk, hv := range r.Header {
		for _, v := range hv {
			w.Header().Add(hk, v)
		}
	}

	w.Write(body)
}

func Test_Server_Base(t *testing.T) {
t.Run(
	"",
	func(t *tesing.T) {

	},
)
}

func Test_Server_DoRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(_Echo_Handler))

	defer func() {
		ts.Close()
	}()

	TEST_DATA := "HELLO"
	req, _ := http.NewRequest("GET", ts.URL, bytes.NewBufferString(TEST_DATA))

	//任意のHeaderを追加。
	req.Header.Add("X-Header1", "test")
	req.Header.Add("X-Header2", "test")
	req.Header.Add("X-Header3", "test, aaaa")

	DoRequest(
		ts.URL,
		req,
		func(response *http.Response) error {
			//HeaderもEchoされているか確認するため、Request-HeaderからKeyを基に値を消去する。
			for k, _ := range response.Header {
				req.Header.Del(k)
			}

			//Request-Headerが空ならEchoされている。
			if len(req.Header) != 0 {
				t.Fatalf("HTTP-Header is not echoed.")
			}

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			if TEST_DATA != string(body) {
				t.Fatalf("Response data are unexpected.")
			}

			return nil
		},
	)
}
