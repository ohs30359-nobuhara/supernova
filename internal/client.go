package internal

import (
	"supernova/pkg"
	"supernova/pkg/slack"
)

type ClientSet struct {
	Slack *slack.Client
}

func NewClientSet(option pkg.SettingsOption) *ClientSet {
	client := &ClientSet{}

	if option.Slack != nil {
		client.Slack = &slack.Client{
			Token:           option.Slack.Token,
			Url:             option.Slack.Url,
			DefaultIcon:     *option.Slack.Icon,
			DefaultUsername: *option.Slack.UserName,
		}
	}
	return client
}
