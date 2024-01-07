package template

import (
	"errors"
	"fmt"
	"os"
	"supernova/internal"
	"supernova/pkg"
	"supernova/pkg/slack"
)

type Runner interface {
	Run() Output
}

type OutputStatus string
type OutputType string

const (
	OutputStatusOK       OutputStatus = "OK"
	OutputStatusWarn     OutputStatus = "WARN"
	OutputStatusDanger   OutputStatus = "DANGER"
	OutputStatusCritical OutputStatus = "CRITICAL"
	OutputTypeText       OutputType   = "TEXT"
	OutputTypeJson       OutputType   = "JSON"
	OutputTypeFile       OutputType   = "FILE"
)

type Output struct {
	Bodies []OutputBody
}

type OutputBody struct {
	Body        []byte
	ContentType OutputType
	Status      OutputStatus
}

func (o *Output) SetBody(body OutputBody) Output {
	o.Bodies = append(o.Bodies, body)
	return *o
}

func (o *Output) Post(client *internal.ClientSet, option *pkg.OutputOption, jobName string) error {
	logger := pkg.GetLogger()

	for _, body := range o.Bodies {
		if option.Slack != nil {
			if e := o.postSlack(client, body, option, jobName); e != nil {
				logger.Error("failed to send slack. " + e.Error())
			}
		}

		if option.File != nil {
			if e := o.putFile(body, option); e != nil {
				logger.Error("failed to put file. " + e.Error())
			}
		}

		if body.ContentType == OutputTypeFile {
			logger.Info(fmt.Sprintf("status: %s, type: %s", body.Status, body.ContentType))
		} else {
			logger.Info(fmt.Sprintf("status: %s, type: %s, body: %s", body.Status, body.ContentType, body.Body))
		}
	}
	return nil
}

// postSlack Slackに結果を送信する
func (o *Output) postSlack(client *internal.ClientSet, body OutputBody, option *pkg.OutputOption, jobName string) error {
	if client.Slack == nil {
		return errors.New("slack is not configured")
	}

	var attachment slack.Attachment

	// TODO: fileの場合を考慮できていないので追加実装する
	switch body.Status {
	case OutputStatusOK:
		attachment.Title = "job executed successfully"
		attachment.Text = fmt.Sprintf("```%s```", string(body.Body))
		attachment.Color = "#36a64f"

	case OutputStatusWarn:
		attachment.Title = "job execution failed"
		attachment.Text = fmt.Sprintf("```%s```", string(body.Body))
		attachment.Color = "#FF0000"
	case OutputStatusDanger:
		attachment.Title = "job execution failed"
		attachment.Text = fmt.Sprintf("```%s```", string(body.Body))
		attachment.Color = "#FF0000"
	}

	msg := slack.Message{
		Text:        fmt.Sprintf("%s supernova step '%s'", option.Slack.Mention, jobName),
		Attachments: []slack.Attachment{attachment},
		Channel:     option.Slack.Channel,
		Link:        true,
	}

	return client.Slack.Post(msg)
}

func (o *Output) putFile(body OutputBody, option *pkg.OutputOption) error {
	if option.File == nil {
		return nil
	}

	// TODO: ここ不具合 (mkdir で作ったdir配下にファイルが生成されていない)
	if e := os.Mkdir(option.File.Dir, os.ModePerm); e != nil {
		return e
	}

	if e := os.WriteFile(option.File.FileName, body.Body, 0644); e != nil {
		return e
	}
	return nil
}
