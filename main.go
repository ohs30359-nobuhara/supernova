package main

import (
	"errors"
	"flag"
	"os"
	"supernova/internal"
	"supernova/pkg"
	"supernova/template"

	"github.com/mitchellh/mapstructure"

	"go.uber.org/zap"
)

func main() {
	confPath := flag.String("config", "./config.yaml", "-config=./config.yaml")
	flag.Parse()

	logger := pkg.GetLogger()
	conf, e := pkg.NewConfig(*confPath)
	if e != nil {
		logger.Error("The pipeline format is incorrect. Please refer to the README.", zap.Error(e))
		os.Exit(1)
	}

	client := internal.NewClientSet(conf.Settings)

	for _, step := range conf.Steps {
		runner, e := getRunner(step.Template)
		if e != nil {
			logger.Error("Failed to load the template configuration", zap.String("name", step.Template), zap.Error(e))
			os.Exit(1)
		}

		if e := execJob(&step, runner, client); e != nil {
			logger.Error("job execution failed", zap.String("name", step.Template), zap.Error(e))
			os.Exit(1)
		}
	}
	os.Exit(0)
}

// getRunner: テンプレートに基づいて対応するランナーを取得する関数
func getRunner(templateName string) (template.Runner, error) {
	switch templateName {
	case "curl":
		return &template.CurlTemplate{}, nil
	case "html":
		return &template.HtmlTemplate{}, nil
	case "shell":
		return &template.ShellTemplate{}, nil
	case "redis":
		return &template.RedisTemplate{}, nil
	default:
		return nil, errors.New("unknown template: " + templateName)
	}
}

// execJob templateの実行
func execJob(step *pkg.StepOption, target template.Runner, client *internal.ClientSet) error {
	if e := mapstructure.Decode(step.Option, target); e != nil {
		return errors.New("template decode error:" + e.Error())
	}
	output := target.Run()
	return output.Post(client, &step.Output, step.Name)
}
