package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"supernova/pkg"
	"supernova/pkg/slack"
	"supernova/template"

	"github.com/mitchellh/mapstructure"

	"go.uber.org/zap"
)

type TemplateRunner interface {
	Run() template.Result
}

func main() {
	confPath := flag.String("config", "./config.yaml", "-config=./config.yaml")
	flag.Parse()

	logger := pkg.GetLogger()
	conf, e := pkg.NewConfig(*confPath)
	if e != nil {
		logger.Error("The pipeline format is incorrect. Please refer to the README.", zap.Error(e))
		panic(e)
	}

	slackClient := slack.Client{
		Token:           conf.Settings.Slack.Token,
		Url:             conf.Settings.Slack.Url,
		DefaultIcon:     *conf.Settings.Slack.Icon,
		DefaultUsername: *conf.Settings.Slack.UserName,
	}

	for _, step := range conf.Steps {
		var err error
		switch step.Template {
		case "curl":
			var t template.CurlTemplate
			err = ExecJob(step, &t, slackClient)
		case "html":
			var t template.HtmlTemplate
			err = ExecJob(step, &t, slackClient)
		case "shell":
			var t template.ShellTemplate
			err = ExecJob(step, &t, slackClient)
		case "redis":
			var t template.RedisTemplate
			err = ExecJob(step, &t, slackClient)
		}

		if err != nil {
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func ExecJob(step pkg.StepOption, target TemplateRunner, client slack.Client) error {
	if e := mapstructure.Decode(step.Option, target); e != nil {
		return errors.New("template decode error:" + e.Error())
	}

	r := target.Run()
	if e := Output(r, step.Output, step.Name, client); e != nil {
		fmt.Println(e)
	}
	return r.Err
}

func Output(result template.Result, option pkg.OutputOption, name string, client slack.Client) error {
	logger := pkg.GetLogger()

	if result.Err != nil {
		logger.Error("The step execution failed.", zap.String("job", name), zap.Error(result.Err))
	} else {
		logger.Info("The step execution success.", zap.String("name", name), zap.String("result", result.Body))
	}

	if option.File != nil {
		if e := os.Mkdir(option.File.Dir, os.ModePerm); e != nil {
			return e
		}

		if e := os.WriteFile(option.File.FileName, []byte(result.Body), 0644); e != nil {
			return e
		}
	}

	if option.Slack != nil {
		var attachment slack.Attachment
		if result.Err != nil {
			attachment.Title = "job executed successfully"
			attachment.Text = fmt.Sprintf("```%s```", result.Body)
			attachment.Color = "#36a64f"
		} else {
			attachment.Title = "job execution failed"
			attachment.Text = fmt.Sprintf("```%s```", result.Err.Error())
			attachment.Color = "#FF0000"
		}

		body := slack.Message{
			Text:        fmt.Sprintf("%s supernova step '%s'", option.Slack.Mention, name),
			Attachments: []slack.Attachment{attachment},
			Channel:     option.Slack.Channel,
			Link:        true,
		}

		if e := client.Post(body); e != nil {
			logger.Error("failed to notify slack.")
		}
	}
	return nil
}
