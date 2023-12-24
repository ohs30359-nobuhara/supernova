package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Text        string       `json:"text"`
	Channel     string       `json:"channel"`
	UserName    *string      `yaml:"username"`
	Icon        *string      `yaml:"icon"`
	Link        bool         `yaml:"link"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text  string `json:"text"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type Client struct {
	Url             string
	Token           string
	DefaultIcon     string
	DefaultUsername string
}

func New(host, token string) Client {
	return Client{
		Url:   host,
		Token: token,
	}
}

func (c Client) Post(msg Message) error {
	if msg.Icon != nil {
		msg.Icon = &c.DefaultIcon
	}

	if msg.UserName != nil {
		msg.UserName = &c.DefaultUsername
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	fmt.Println("Message sent successfully!")
	return nil
}
