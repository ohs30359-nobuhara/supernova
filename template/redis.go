package template

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

type RedisTemplate struct {
	Cluster  *[]string
	Single   *string
	Password string
	Commands []RedisCommand
}

type RedisCommand struct {
	Action    string
	Key       string
	Set       *string
	expireMin *int
}

// Run templateの実行
func (t RedisTemplate) Run() error {
	client, err := t.createRedisInstance()
	if err != nil {
		return err
	}
	defer client.Close()

	// 先に接続確認
	if err := client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis. %w", err)
	}

	for _, cmd := range t.Commands {
		switch strings.ToUpper(cmd.Action) {
		case "GET":
			result, err := client.Get(context.Background(), cmd.Key).Result()
			if err != nil {
				return fmt.Errorf("GET command has failed. %w", err)
			}
			fmt.Println(result)
		case "SET":
			expire := time.Duration(*cmd.expireMin) * time.Minute
			result, err := client.Set(context.Background(), cmd.Key, *cmd.Set, expire).Result()
			if err != nil {
				return fmt.Errorf("SET command has failed. %w", err)
			}
			fmt.Println(result)
		case "DELETE":
			result, err := client.Del(context.Background(), cmd.Key).Result()
			if err != nil {
				return fmt.Errorf("DELETE command has failed. %w", err)
			}
			fmt.Println(result)
		}
	}

	return nil
}

// createRedisInstance RedisのInstanceを生成
func (t RedisTemplate) createRedisInstance() (redis.UniversalClient, error) {
	if t.Cluster != nil {
		options := &redis.ClusterOptions{
			Addrs:    *t.Cluster,
			Password: t.Password,
		}
		client := redis.NewClusterClient(options)
		return client, nil
	}

	if t.Single != nil {
		options := &redis.Options{
			Addr:     *t.Single,
			Password: t.Password,
		}
		client := redis.NewClient(options)
		return client, nil
	}
	return nil, errors.New("there is an error in the Redis connection information")
}
