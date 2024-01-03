package infra

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/pharmaniaga/auth-user/domain/model/general"
	"github.com/spf13/viper"
)

var (
	ctx = context.Background()
)

func NewRedisClient(cfg general.RedisAccount) (*redis.Client, error) {

	options := &redis.Options{
		Addr:       fmt.Sprintf("%s:%d", viper.GetString("REDIS.URL"), viper.GetInt("REDIS.PORT")),
		Password:   viper.GetString("REDIS.PASSWORD"),
		Username:   viper.GetString("REDIS.USERNAME"),
		DB:         0,
		MaxRetries: 3,
		Network:    "tcp",
		TLSConfig:  &tls.Config{},
	}

	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis : %w", err)
	}

	return client, nil
}
