package pkg

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Job定義
	Steps []StepOption `yaml:"steps"`
	// 共通設定
	Settings SettingsOption `yaml:"settings"`
}

type StepOption struct {
	// Job名
	Name string `yaml:"name"`
	// 利用Template
	Template string `yaml:"template"`
	// Templateごとの独自設定
	Option map[string]interface{} `yaml:"option"`
	// 出力設定 (標準出力は設定に関わらず常に出力)
	Output OutputOption `yaml:"output"`
}

type OutputOption struct {
	// Slackへの通知設定
	Slack *OutputSlackOption `yaml:"slack"`
	// ファイルの出力設定
	File *OutputFileOption `yaml:"file"`
}

type OutputSlackOption struct {
	// 送信先チャンネル
	Channel string `yaml:"channel"`
	// 送信時のメンション
	Mention string `yaml:"mention"`
	// 出力タイプ (text or file)
	Type string `yaml:"type"`
}

type OutputFileOption struct {
	// 出力ファイル名
	FileName string `yaml:"fileName"`
	// 出力先ディレクトリ
	Dir string `yaml:"dir"`
}

type SettingsOption struct {
	// Slack共通設定
	Slack *SettingsSlackOption `yaml:"slack"`
}

type SettingsSlackOption struct {
	// Slack API URL
	Url string `yaml:"url"`
	// 認証トークン
	Token string `yaml:"token"`
	// 表示アイコン
	Icon *string `yaml:"icon"`
	// 通知ユーザー名 (表示)
	UserName *string `yaml:"userName"`
}

func NewConfig(path string) (Config, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return Config{}, e
	}

	var config Config
	if e := yaml.Unmarshal(b, &config); e != nil {
		return Config{}, e
	}
	return config, nil
}
