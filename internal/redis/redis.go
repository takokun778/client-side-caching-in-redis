package redis

import (
	"context"

	"github.com/redis/rueidis"
)

type Client struct {
	rueidis.Client
}

func New() *Client {
	ctx := context.Background()

	cli, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
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
