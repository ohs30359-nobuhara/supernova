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
	Diff   *struct {
		Url     string `yaml:"url"`
		WaitSec int    `yaml:"waitSec"`
		Slack   string `yaml:"string"`
	} `yaml:"diff"`
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
func (t HtmlTemplate) Run() Result {
	logger := pkg.GetLogger()
	command := "node ./browser/dist/main.js "

	// スクリーンショット処理
	if t.Screenshot != nil {
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
