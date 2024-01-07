package template

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
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
		Text   *string `yaml:"text"`
		Api    *string `yaml:"api"`
		File   *string `yaml:"file"`
	} `yaml:"expect"`
}

// Run templateの実行
func (t CurlTemplate) Run() Output {
	var output Output

	status, body, e := t.request()
	if e != nil {
		return output.SetBody(OutputBody{
			Body:        []byte("request failed. " + e.Error()),
			Status:      OutputStatusDanger,
			ContentType: OutputTypeText,
		})
	}

	if e := t.verifyResponse(status, body); e != nil {
		return output.SetBody(OutputBody{
			Body:        []byte("validation result was incorrect. " + e.Error()),
			Status:      OutputStatusWarn,
			ContentType: OutputTypeText,
		})
	}

	return output.SetBody(OutputBody{
		Body:        body,
		Status:      OutputStatusOK,
		ContentType: OutputTypeText,
	})
}

// request HTTP Requestを投げる
func (t CurlTemplate) request() (int, []byte, error) {
	var requestBody io.Reader
	if t.Body != nil {
		requestBody = bytes.NewBuffer([]byte(*t.Body))
	}

	req, err := http.NewRequest(t.Method, t.URL, requestBody)
	if err != nil {
		return 0, nil, err
	}

	for _, h := range t.Header {
		parts := strings.Split(h, ":")
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	req.Header.Set("Cookie", t.Cookie)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, respBody, nil
}

// verifyResponse レスポンスを検証する
func (t CurlTemplate) verifyResponse(status int, body []byte) error {
	if err := t.compareStatus(status); err != nil {
		return err
	}

	if err := t.compareText(body); err != nil {
		return err
	}

	if err := t.compareApi(body, status); err != nil {
		return err
	}

	if err := t.compareFile(body); err != nil {
		return err
	}

	return nil
}

// compareStatus ステータスを比較
func (t CurlTemplate) compareStatus(status int) error {
	if t.Expect.Status != nil && status != *t.Expect.Status {
		return errors.New("response status code did not match")
	}
	return nil
}

// compareText テキストを比較
func (t CurlTemplate) compareText(body []byte) error {
	if t.Expect.Text != nil && *t.Expect.Text != string(body) {
		return errors.New("response body did not match")
	}
	return nil
}

// compareApi APIを叩いて比較
func (t CurlTemplate) compareApi(body []byte, status int) error {
	if t.Expect.Api != nil {
		expectStatus, expectBody, err := t.request()
		if err != nil {
			return errors.New("The request to the comparison API has failed. " + err.Error())
		}
		if expectStatus != status {
			return errors.New("response status code did not match")
		}

		if !bytes.Equal(body, expectBody) {
			return errors.New("response body did not match")
		}
	}
	return nil
}

// compareFile ファイルを比較
func (t CurlTemplate) compareFile(body []byte) error {
	if t.Expect.File != nil {
		buf, err := os.ReadFile(*t.Expect.File)
		if err != nil {
			return errors.New("file not found: " + err.Error())
		}
		if !bytes.Equal(body, buf) {
			return errors.New("response body did not match")
		}
	}
	return nil
}
