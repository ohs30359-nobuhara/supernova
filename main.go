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
		panic(e)
	}

	client := internal.NewClientSet(conf.Settings)

	for _, step := range conf.Steps {
		var err error
		switch step.Template {
		case "curl":
			var t template.CurlTemplate
			err = execJob(&step, &t, client)
		case "html":
			var t template.HtmlTemplate
			err = execJob(&step, &t, client)
		case "shell":
			var t template.ShellTemplate
			err = execJob(&step, &t, client)
		case "redis":
			var t template.RedisTemplate
			err = execJob(&step, &t, client)
		}

		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}
	os.Exit(0)
}

// execJob templateの実行
func execJob(step *pkg.StepOption, target template.Runner, client *internal.ClientSet) error {
	if e := mapstructure.Decode(step.Option, target); e != nil {
		return errors.New("template decode error:" + e.Error())
	}
	output := target.Run()
	return output.Post(client, &step.Output, step.Name)
}
