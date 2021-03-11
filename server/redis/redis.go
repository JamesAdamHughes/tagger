package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"reflect"
	"time"
)

var redisClient *redis.Client

func init() {
	fmt.Printf("init connection to redis....\n")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.TODO()

	fmt.Printf("\n%+v\n", redisClient.Ping(ctx))
}

func Set(key string, value interface{}, ttl int) error {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return errors.New("Set requires a pointer to store an entity")
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(value)
	if err != nil {
		return err
	}

	ctx := context.TODO()
	if err = redisClient.Set(ctx, key, buffer.Bytes(), time.Duration(ttl)*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func Get(key string, i interface{}) error {
	if reflect.ValueOf(i).Kind() != reflect.Ptr {
		return errors.New("Set requires a pointer to store an entity")
	}

	ctx := context.TODO()
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Printf("\nKey: %s does not exist \n", key)
		return nil
	} else if err != nil {
		return err
	}

	var buffer = bytes.NewBuffer([]byte(val))
	dec := gob.NewDecoder(buffer)
	if err = dec.Decode(i); err != nil {
		return err
	}

	return nil
}
