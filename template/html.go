package template

import (
	"fmt"
	diffimage "github.com/murooka/go-diff-image"
	"supernova/pkg/browser"
	"supernova/pkg/img"
)

type HtmlTemplate struct {
	URL    string   `yaml:"url"`
	Header []string `yaml:"header"`
	Cookie string   `yaml:"string"`
	Diff   *struct {
		Url     string `yaml:"url"`
		WaitSec int    `yaml:"waitSec"`
		Slack   string `yaml:"string"`
	} `yaml:"diff"`
	Screenshot *struct {
		WaitSec int    `yaml:"waitSec"`
		Slack   string `yaml:"string"`
	} `yaml:"screenshot"`
	/*
		CoreWebVital *struct {
			WaitSec int    `yaml:"waitSec"`
			Slack   string `yaml:"string"`
		} `yaml:"coreWebVital"`
	*/
}

// Run Templateの実行
func (t HtmlTemplate) Run() error {
	page := browser.NewPage(t.URL, t.Header, t.Cookie)

	// スクリーンショット処理
	if t.Screenshot != nil {
		if e := t.screenshots(page); e != nil {
			return e
		}
	}

	if t.Diff != nil {
		if e := t.diff(page); e != nil {
			return e
		}
	}

	/* APIを使えないので一旦無効化
	if t.CoreWebVital != nil {
		if e := coreWebVital(page, t.CoreWebVital.WaitSec); e != nil {
			return e
		}
	}
	*/

	return nil
}

// screenshots スクリーンショットを撮る
func (t HtmlTemplate) screenshots(page browser.Page) error {
	screenshotImg, e := page.Screenshot(t.Screenshot.WaitSec)
	if e != nil {
		return e
	}
	if e := img.SaveImageToFile(*screenshotImg, "./screenshots.png"); e != nil {
		return e
	}
	return nil
}

// diff ページとの差分を取得する
func (t HtmlTemplate) diff(page browser.Page) error {
	actual, e := page.Screenshot(t.Diff.WaitSec)
	if e != nil {
		return e
	}

	diffPage := browser.NewPage(t.Diff.Url, t.Header, t.Cookie)
	expect, e := diffPage.Screenshot(t.Diff.WaitSec)
	if e != nil {
		return e
	}

	diff := diffimage.DiffImage(*actual, *expect)
	return img.SaveImageToFile(diff, "diff.png")
}

func coreWebVital(browser browser.Page, sec int) error {
	vital, e := browser.GetCoreWebVital(sec)
	if e != nil {
		return e
	}
	fmt.Println(vital)
	return nil
}
