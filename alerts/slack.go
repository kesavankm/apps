package alerts

import (
	"context"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	ctx      context.Context
	sclient  *slack.Client
	apiToken string
}

func NewSlackClient(ctx context.Context, apiToken string) (client *SlackClient, err error) {
	// log.Printf("[Slack-New] Enter")
	c := &SlackClient{ctx: ctx, apiToken: apiToken}
	c.sclient = slack.New(apiToken)
	return c, nil
}

func (c *SlackClient) PostMessage(channelID, message string) (err error) {
	// log.Printf("[Slack-PostMessage] Enter")
	rchannelID, ts, err := c.sclient.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		panic(err)
	}
	// log.Printf("[Slack-PostMessage] Sent text message %s to Channel %s @ %s\n", message, rchannelID, ts)
	return nil
}
