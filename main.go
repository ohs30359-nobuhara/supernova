package main

import (
	"errors"
	"os"
	"supernova/pkg"
	"supernova/template"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
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
		logger.Error("The pipeline format is incorrect. Please refer to the README.", zap.Error(e))
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
		case "shell":
			var t template.ShellTemplate
			templateError = ExecJob(step, &t)
		case "redis":
			var t template.RedisTemplate
			templateError = ExecJob(step, &t)
		}

		if templateError != nil {
			logger.Error("The step execution failed.", zap.String("name", step.Name), zap.Error(templateError))
			os.Exit(1)
		} else {
			logger.Info("The step execution Success", zap.String("name", step.Name))
		}
	}
	os.Exit(0)
}
