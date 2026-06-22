package redis

import (
	"backend/internal/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	connCli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	err := connCli.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	if err = redisotel.InstrumentTracing(connCli); err != nil {
		return nil, err
	}

	return connCli, nil
}
