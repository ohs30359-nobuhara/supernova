package main

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"supernova/pkg"
	"supernova/template"
)

type TemplateRunner interface {
	Run()
}

func ExecJob(step pkg.Step, target TemplateRunner) error {
	if e := mapstructure.Decode(step.Option, target); e != nil {
		return errors.New("template decode error:" + e.Error())
	}
	target.Run()
	return nil
}

func main() {
	cnf, e := pkg.NewConfig("./config.yaml")
	if e != nil {
		panic(e)
	}

	for _, step := range cnf.Steps {
		switch step.Template {
		case "curl":
			var t template.CurlTemplate
			ExecJob(step, &t)
		case "html":
			var t template.HtmlTemplate
			ExecJob(step, &t)
		}
	}
}
