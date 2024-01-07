package template

import (
	"context"
	"errors"
	"strings"
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
func (t RedisTemplate) Run() Output {
	var output Output
	client, e := t.createRedisInstance()
	if e != nil {
		return output.SetBody(OutputBody{
			Body:        []byte("failed create redis instance. " + e.Error()),
			ContentType: OutputTypeText,
			Status:      OutputStatusDanger,
		})
	}
	defer client.Close()

	// 先に接続確認
	if e := client.Ping(context.Background()).Err(); e != nil {
		return output.SetBody(OutputBody{
			Body:        []byte("failed to connect to Redis. " + e.Error()),
			ContentType: OutputTypeText,
			Status:      OutputStatusDanger,
		})
	}

	for _, cmd := range t.Commands {
		switch strings.ToUpper(cmd.Action) {
		case "GET":
			if result, e := client.Get(context.Background(), cmd.Key).Result(); e != nil {
				return output.SetBody(OutputBody{
					Body:        []byte("GET command has failed. " + e.Error()),
					ContentType: OutputTypeText,
					Status:      OutputStatusDanger,
				})
			} else {
				output.SetBody(OutputBody{
					Body:        []byte(result),
					ContentType: OutputTypeText,
					Status:      OutputStatusOK,
				})
			}

		case "SET":
			var expire time.Duration
			if cmd.expireMin != nil {
				expire = time.Duration(*cmd.expireMin) * time.Minute
			} else {
				expire = 0
			}
			if _, e := client.Set(context.Background(), cmd.Key, *cmd.Value, expire).Result(); e != nil {
				return output.SetBody(OutputBody{
					Body:        []byte("SET command has failed. " + e.Error()),
					ContentType: OutputTypeText,
					Status:      OutputStatusDanger,
				})
			} else {
				output.SetBody(OutputBody{
					Body:        []byte("SET command was successful."),
					ContentType: OutputTypeText,
					Status:      OutputStatusOK,
				})
			}
		case "DELETE":
			if _, e := client.Del(context.Background(), cmd.Key).Result(); e != nil {
				return output.SetBody(OutputBody{
					Body:        []byte("DELETE command has failed. " + e.Error()),
					ContentType: OutputTypeText,
					Status:      OutputStatusDanger,
				})
			} else {
				output.SetBody(OutputBody{
					Body:        []byte("DELETE command was successful."),
					ContentType: OutputTypeText,
					Status:      OutputStatusOK,
				})
			}
		}
	}
	return output
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
