package template

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

type CurlTemplate struct {
	URL    string   `yaml:"url"`
	Method string   `yaml:"method"`
	Body   *string  `yaml:"body"`
	Header []string `yaml:"header"`
	Cookie string   `yaml:"cookie"`
	Expect struct {
		Status *int    `yaml:"status"`
		Body   *string `yaml:"body"`
		Equal  *string `yaml:"equal"`
	} `yaml:"expect"`
}

// Run templateの実行
func (t CurlTemplate) Run() error {
	resp, e := request(t.Method, t.URL, t.Cookie, t.Header, t.Body)
	if e != nil {
		return e
	}
	
	if t.Expect.Status != nil {
		if resp.StatusCode != *t.Expect.Status {
			return errors.New("status codeが一致しませんでした")
		}
	}

	if t.Expect.Body != nil {
		body, _ := io.ReadAll(resp.Body)
		if *t.Expect.Body != string(body) {
			return errors.New("bodyが一致しませんでした")
		}
	}

	if t.Expect.Equal != nil {
		equalResp, e := request(t.Method, *t.Expect.Equal, t.Cookie, t.Header, t.Body)
		if e != nil {
			return errors.New("比較側のレスポンス取得に失敗しました: " + e.Error())
		}
		if equalResp.StatusCode != resp.StatusCode {
			return errors.New("status codeが一致しませんでした")
		}
		// bodyの比較をする
		equalBody, _ := io.ReadAll(equalResp.Body)
		originalBody, _ := io.ReadAll(resp.Body)
		if string(equalBody) != string(originalBody) {
			return errors.New("bodyが一致しませんでした")
		}
	}

	return nil
}

// request HTTP Requestを投げる
func request(method, url, cookie string, header []string, bodyStr *string) (*http.Response, error) {
	var body io.Reader
	if bodyStr != nil {
		body = bytes.NewBuffer([]byte(*bodyStr))
	}

	req, e := http.NewRequest(method, url, body)
	if e != nil {
		return nil, e
	}

	for _, h := range header {
		parts := strings.Split(h, ":")
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	req.Header.Set("Cookie", cookie)

	client := http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		return nil, e
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	return resp, e
}
