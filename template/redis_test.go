package template

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedisTemplate_Run(t *testing.T) {
	// モックのRedisサーバーを起動
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	addr := "localhost:6379"

	// テスト用のRedisTemplateインスタンスを作成
	template := RedisTemplate{
		Single:   &addr,
		Password: "",
		Commands: []RedisCommand{
			{Action: "SET", Key: "test_key", Value: stringPointer("test_value"), expireMin: intPointer(5)},
			{Action: "GET", Key: "test_key"},
			{Action: "DELETE", Key: "test_key"},
		},
	}

	// テストを実行
	err = template.Run()
	assert.NoError(t, err)

	// テスト用のRedisTemplateインスタンスを作成（エラーケース）
	templateWithInvalidConnection := RedisTemplate{
		Cluster:  nil,
		Single:   nil,
		Password: "",
		Commands: []RedisCommand{},
	}

	// エラーケースのテストを実行
	err = templateWithInvalidConnection.Run()
	assert.Error(t, err)
	assert.Equal(t, "there is an error in the Redis connection information", err.Error())
}

func stringPointer(s string) *string {
	return &s
}

func intPointer(i int) *int {
	return &i
}
