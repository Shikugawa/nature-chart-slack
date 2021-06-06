package slack

import (
	"github.com/slack-go/slack"
)

type SlackClient struct {
	client *slack.Client
}

func NewSlackClient(atoken string) *SlackClient {
	return &SlackClient{
		client: slack.New(atoken),
	}
}

func (s *SlackClient) Post(dstPath string, dstChannel string) error {
	params := slack.FileUploadParameters{
		Title:    "Current Temp",
		File:     dstPath,
		Channels: []string{dstChannel},
	}
	_, err := s.client.UploadFile(params)
	if err != nil {
		return err
	}

	return nil
}
