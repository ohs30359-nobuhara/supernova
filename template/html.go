package template

import "fmt"

type HtmlTemplate struct {
	URL          string `yaml:"url"`
	Snapshot     string `yaml:"snapshot"`
	CoreWebVital string `yaml:"coreWebVital"`
}

func (t HtmlTemplate) Run() {
	fmt.Println("call html")
}
