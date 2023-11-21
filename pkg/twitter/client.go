package twitter

import (
	"context"
	"zen/pkg/logger"
)

type Client struct {
	ctx               context.Context
	log               *logger.Logger
	twitterCredential TwitterCredential
}

type TwitterCredential struct {
	CT0       string
	AuthToken string
}

func NewClient(ctx context.Context, log *logger.Logger, twitterCredential TwitterCredential) *Client {
	return &Client{
		ctx:               ctx,
		log:               log.WithField("component", "twitter_client"),
		twitterCredential: twitterCredential,
	}
}
