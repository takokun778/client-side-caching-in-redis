package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/takokun778/client-side-caching-in-redis/internal/redis"
)

func main() {
	ctx := context.Background()

	client := redis.New()

	defer client.Close()

	{
		cmd := client.B().Set().Key("key").Value("value").Build()

		if err := client.Do(ctx, cmd).Error(); err != nil {
			panic(err)
		}
	}

	{
		cmd := client.B().Get().Key("key").Build()

		got, err := client.Do(ctx, cmd).ToString()
		if err != nil {
			panic(err)
		}

		log.Printf("got: %s", got)
	}

	{
		cmd := client.B().Get().Key("key").Cache()

		got, err := client.DoCache(ctx, cmd, time.Second).ToString()
		if err != nil {
			panic(err)
		}

		log.Printf("got: %s", got)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		fmt.Println(i)
	}

	{
		cmd := client.B().Set().Key("key").Value("valuevalue").Build()

		if err := client.Do(ctx, cmd).Error(); err != nil {
			panic(err)
		}
	}

	// {
	// 	cmd := client.B().Del().Key("key").Build()

	// 	if err := client.Do(ctx, cmd).Error(); err != nil {
	// 		panic(err)
	// 	}
	// }
}
