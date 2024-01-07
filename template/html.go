package template

import (
	"os/exec"
	"supernova/pkg"

	"github.com/mattn/go-shellwords"
)

type HtmlTemplate struct {
	URL    string   `yaml:"url"`
	Header []string `yaml:"header"`
	Cookie string   `yaml:"string"`
	// 検証
	Expect struct {
		// 画像による差分比較の有効化
		Diff *struct {
			// 比較対象のURL
			Url string `yaml:"url"`
			// 撮影までの待機時間
			WaitSec int `yaml:"waitSec"`
		} `yaml:"diff"`
	} `yaml:"expect"`
	// スクリーンショットの取得
	Screenshot bool `yaml:"screenshot"`
	// CoreWebVitalの取得
	CoreWebVital *struct {
		// 出力形式 html or json
		Format string `yaml:"format"`
	} `yaml:"coreWebVital"`
}

// Run Templateの実行
func (t HtmlTemplate) Run() Result {
	logger := pkg.GetLogger()
	command := "node ./browser/dist/main.js "

	// スクリーンショット処理
	if t.Screenshot {
		command += "--screenshot screenshot.png "
	}

	// APIを使えないので一旦無効化
	if t.CoreWebVital != nil {
		command += "--performance true "
	}

	command += t.URL

	logger.Info(command)
	args, e := shellwords.Parse(command)
	if e != nil {
		return NewResultError("", DANGER, e)
	}
	output, e := exec.Command(args[0], args[1:]...).CombinedOutput()
	if e != nil {
		return NewResultError("", DANGER, e)
	}

	return NewResultSuccess(string(output))
}
