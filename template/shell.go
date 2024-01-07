package template

import (
	"os"
	"os/exec"
)

type ShellTemplate struct {
	Script string  `yaml:"script"` // 実行スクリプト
	Dir    *string `yaml:"dir"`    // 実行ディレクトリの指定 (未指定ならカレントディレクトリ)
}

func (t ShellTemplate) Run() Output {
	var output Output

	out, e := t.execScript()
	if e != nil {
		return output.SetBody(OutputBody{
			Body:        []byte("failed exec command. " + e.Error()),
			Status:      OutputStatusDanger,
			ContentType: OutputTypeText,
		})
	}

	return output.SetBody(OutputBody{
		Body:        out,
		Status:      OutputStatusOK,
		ContentType: OutputTypeText,
		Metadata:    map[OutputMetadata]string{OutputMetadataFileName: "shell.txt"},
	})
}

// execScript shell commandを実行
func (t ShellTemplate) execScript() ([]byte, error) {
	path := "temp.sh"
	temp, e := os.CreateTemp(".", path)
	if e != nil {
		return nil, e
	}

	if _, e := temp.Write([]byte(t.Script)); e != nil {
		return nil, e
	}

	defer os.Remove(temp.Name())

	cmd := exec.Command("sh", "./"+temp.Name())

	if t.Dir != nil {
		cmd.Dir = *t.Dir
	}

	output, e := cmd.CombinedOutput()
	if e != nil {
		return nil, e
	}
	return output, nil
}
