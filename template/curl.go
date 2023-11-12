package template

import (
	"net/http"
	"strings"
)

type CurlTemplate struct {
	URL    string   `yaml:"url"`
	Method string   `yaml:"method"`
	Body   *string  `yaml:"body"`
	Header []string `yaml:"header"`
	Cookie string   `yaml:"cookie"`
}

func (t CurlTemplate) Run() error {
	req, e := http.NewRequest(t.Method, t.URL, nil)
	if e != nil {
		return e
	}

	for _, h := range t.Header {
		parts := strings.Split(h, ":")
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	req.Header.Set("Cookie", t.Cookie)

	client := http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	return nil
}
