package template

import (
	"fmt"
	"supernova/pkg/browser"
	"supernova/pkg/image"
)

type HtmlTemplate struct {
	URL        string   `yaml:"url"`
	Header     []string `yaml:"header"`
	Cookie     string   `yaml:"string"`
	Screenshot *struct {
		WaitSec int    `yaml:"waitSec"`
		Slack   string `yaml:"string"`
	} `yaml:"screenshot"`
	CoreWebVital *struct {
		WaitSec int    `yaml:"waitSec"`
		Slack   string `yaml:"string"`
	} `yaml:"coreWebVital"`
}

// Run Templateの実行
func (t HtmlTemplate) Run() error {
	browser := browser.NewHeadlessBrowser(t.URL, t.Header, t.Cookie)

	// スクリーンショット処理
	if t.Screenshot != nil {
		if e := screenshot(browser, t.Screenshot.WaitSec); e != nil {
			return e
		}
	}

	if t.CoreWebVital != nil {
		if e := coreWebVital(browser, t.CoreWebVital.WaitSec); e != nil {
			return e
		}
	}

	return nil
}

// screenshot
func screenshot(browser browser.HeadlessBrowser, sec int) error {
	buf, e := browser.Screenshot(sec)
	if e != nil {
		return e
	}

	img, err := image.ByteArrayToImage(*buf)
	if err != nil {
		return err
	}

	// 画像をファイルに保存する例
	err = image.SaveImageToFile(img, "./screenshot.png")
	if err != nil {
		return err
	}

	return nil
}

func coreWebVital(browser browser.HeadlessBrowser, sec int) error {
	vital, e := browser.GetCoreWebVital(sec)
	if e != nil {
		return e
	}
	fmt.Println(vital)
	return nil
}
