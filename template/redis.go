package template

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"supernova/pkg"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisTemplate struct {
	Cluster  *[]string      `yaml:"cluster"`
	Single   *string        `yaml:"single"`
	Password string         `yaml:"password"`
	Commands []RedisCommand `yaml:"commands"`
}

type RedisCommand struct {
	Action    string  `yaml:"action"`
	Key       string  `yaml:"key"`
	Value     *string `yaml:"value"`
	expireMin *int    `yaml:"expireMin"`
}

// Run templateの実行
func (t RedisTemplate) Run() error {
	logger := pkg.GetLogger()
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
			logger.Info(result)

		case "SET":
			var expire time.Duration
			if cmd.expireMin != nil {
				expire = time.Duration(*cmd.expireMin) * time.Minute
			} else {
				expire = 0
			}
			result, err := client.Set(context.Background(), cmd.Key, *cmd.Value, expire).Result()
			if err != nil {
				return fmt.Errorf("SET command has failed. %w", err)
			}
			logger.Info(result)
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
