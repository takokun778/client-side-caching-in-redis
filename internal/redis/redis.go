package redis

import (
	"context"

	"github.com/redis/rueidis"
)

type Client struct {
	rueidis.Client
}

func New(url string) *Client {
	ctx := context.Background()

	option := rueidis.ClientOption{
		InitAddress: []string{url},
	}

	cli, err := rueidis.NewClient(option)
	if err != nil {
		panic(err)
	}

	if err := cli.Do(ctx, cli.B().Ping().Build()).Error(); err != nil {
		panic(err)
	}

	return &Client{
		Client: cli,
	}
}
