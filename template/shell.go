package template

import (
	"fmt"
	"os"
	"os/exec"
)

type ShellTemplate struct {
	Script string  `yaml:"script"` // 実行スクリプト
	Dir    *string `yaml:"dir"`    // 実行ディレクトリの指定 (未指定ならカレントディレクトリ)
	Output bool    `yaml:"output"` // ログ出力するか
}

func (t ShellTemplate) Run() Result {
	out, e := t.execScript()
	if e != nil {
		return NewResultError("failed exec command,", DANGER, e)
	}

	if t.Output {
		fmt.Println(out)
	}
	return NewResultSuccess("")
}

func (t ShellTemplate) execScript() (string, error) {
	path := "temp.sh"
	temp, e := os.CreateTemp(".", path)
	if e != nil {
		return "", e
	}

	if _, e := temp.Write([]byte(t.Script)); e != nil {
		return "", e
	}

	defer os.Remove(temp.Name())

	cmd := exec.Command("sh", "./"+temp.Name())

	if t.Dir != nil {
		cmd.Dir = *t.Dir
	}

	output, e := cmd.CombinedOutput()
	if e != nil {
		return "", e
	}
	return string(output), nil
}
