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
func (t RedisTemplate) Run() Result {
	logger := pkg.GetLogger()
	client, e := t.createRedisInstance()
	if e != nil {
		return NewResultError("failed create redis instance", DANGER, e)
	}
	defer client.Close()

	// 先に接続確認
	if e := client.Ping(context.Background()).Err(); e != nil {
		return NewResultError("failed to connect to Redis", DANGER, e)
	}

	for _, cmd := range t.Commands {
		switch strings.ToUpper(cmd.Action) {
		case "GET":
			result, err := client.Get(context.Background(), cmd.Key).Result()
			if err != nil {
				return NewResultError("GET command has failed.", DANGER, e)
			}
			logger.Info(result)

		case "SET":
			var expire time.Duration
			if cmd.expireMin != nil {
				expire = time.Duration(*cmd.expireMin) * time.Minute
			} else {
				expire = 0
			}
			result, e := client.Set(context.Background(), cmd.Key, *cmd.Value, expire).Result()
			if e != nil {
				return NewResultError("SET command has failed.", DANGER, e)
			}
			logger.Info(result)
		case "DELETE":
			result, err := client.Del(context.Background(), cmd.Key).Result()
			if err != nil {
				return NewResultError("DELETE command has failed.", DANGER, e)
			}
			fmt.Println(result)
		}
	}
	return NewResultSuccess("")
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
