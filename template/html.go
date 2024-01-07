package template

import (
	"os"
	"os/exec"

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
func (t HtmlTemplate) Run() Output {
	var output Output

	_, e := t.callBrowserCli()
	if e != nil {
		return output.SetBody(OutputBody{
			Body:        []byte("an unexpected error has occurred. Please check the log."),
			Status:      OutputStatusDanger,
			ContentType: OutputTypeJson,
		})
	}

	if t.Screenshot {
		if buf, e := os.ReadFile("./screenshot.png"); e != nil {
			output.SetBody(OutputBody{
				Body:        []byte("failed to take snapshot"),
				Status:      OutputStatusDanger,
				ContentType: OutputTypeFile,
			})
		} else {
			output.SetBody(OutputBody{
				Body:        buf,
				Status:      OutputStatusOK,
				ContentType: OutputTypeFile,
				Metadata:    map[OutputMetadata]string{OutputMetadataFileName: "screenshot.png"},
			})
		}
	}

	if t.CoreWebVital != nil {
		target := ""
		switch t.CoreWebVital.Format {
		case "html":
			target = "lighthouse.html"
		case "json":
			target = "lighthouse.json"
		default:
			return output.SetBody(OutputBody{
				Body:        []byte("core web vital option 'format' must be 'html' or 'json'."),
				Status:      OutputStatusDanger,
				ContentType: OutputTypeFile,
			})
		}

		if buf, e := os.ReadFile(target); e != nil {
			output.SetBody(OutputBody{
				Body:        []byte("failed to take core web vital"),
				Status:      OutputStatusDanger,
				ContentType: OutputTypeFile,
			})
		} else {
			output.SetBody(OutputBody{
				Body:        buf,
				Status:      OutputStatusOK,
				ContentType: OutputTypeFile,
				Metadata:    map[OutputMetadata]string{OutputMetadataFileName: target},
			})
		}
	}
	return output
}

func (t HtmlTemplate) callBrowserCli() ([]byte, error) {
	command := "node ./browser/dist/main.js "

	// スクリーンショット処理
	if t.Screenshot {
		command += "--screenshot screenshot.png "
	}

	// APIを使えないので一旦無効化
	if t.CoreWebVital != nil {
		command += "--performance lighthouse "
	}

	command += t.URL

	args, e := shellwords.Parse(command)
	if e != nil {
		return nil, e
	}

	return exec.Command(args[0], args[1:]...).CombinedOutput()
}
