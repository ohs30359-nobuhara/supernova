package main

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"os"
	"supernova/pkg"
	"supernova/template"
)

type TemplateRunner interface {
	Run() error
}

func ExecJob(step pkg.Step, target TemplateRunner) error {
	if e := mapstructure.Decode(step.Option, target); e != nil {
		return errors.New("template decode error:" + e.Error())
	}
	return target.Run()
}

func main() {
	logger := pkg.GetLogger()
	cnf, e := pkg.NewConfig("./config.yaml")
	if e != nil {
		logger.Error("pipelineの書式が間違っています", zap.Error(e))
		panic(e)
	}

	for _, step := range cnf.Steps {
		var templateError error
		switch step.Template {
		case "curl":
			var t template.CurlTemplate
			templateError = ExecJob(step, &t)
		case "html":
			var t template.HtmlTemplate
			templateError = ExecJob(step, &t)
		}

		if templateError != nil {
			logger.Error("step失敗", zap.String("name", step.Name), zap.Error(templateError))
			os.Exit(1)
		} else {
			logger.Info("成功", zap.String("name", step.Name))
		}
	}
}
