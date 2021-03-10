package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func init() {
	fmt.Printf("init connection to redis....\n")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.TODO()

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
}

func Test() {

}
