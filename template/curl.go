package template

import "fmt"

type CurlTemplate struct {
	URL    string   `yaml:"url"`
	Method string   `yaml:"method"`
	Body   *string  `yaml:"body"`
	Header []string `yaml:"header"`
	Cookie string   `yaml:"cookie"`
}

func (t CurlTemplate) Run() {
	fmt.Println("aaa")
}
