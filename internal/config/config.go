package config

import (
	"backend/internal/log"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type PGConfig struct {
	PGName       string
	PGUser       string
	PGPassword   string
	PGHost       string
	PGPort       int
	MaxPool      int32
	PGTimeout    time.Duration
	ConnAttempts int
}

type RedisConfig struct {
	RedisHost     string
	RedisPort     int
	RedisPassword string
}

type Config struct {
	PGConfig    *PGConfig
	RedisConfig *RedisConfig

	StorageType string
	CacheType   string

	GenerateAttempts int

	CacheTTL time.Duration

	ServiceHost      string
	ServicePort      string
	ShortLinkAddress string
}

const (
	storageType = "STORAGE_TYPE"
	cacheType   = "CACHE_TYPE"

	generateAttempts = "GENERATE_ATTEMPTS"

	PGName         = "PG_NAME"
	PGUser         = "PG_USER"
	PGPassword     = "PG_PASSWORD"
	PGHost         = "PG_HOST"
	PGPort         = "PG_PORT"
	PGTimeout      = "PG_TIMEOUT"
	PGMaxPool      = "PG_MAX_POOL"
	PGConnAttempts = "PG_CONN_ATTEMPTS"

	RedisHost     = "REDIS_HOST"
	RedisPort     = "REDIS_PORT"
	RedisPassword = "REDIS_PASSWORD"
	cacheTTL      = "CACHE_TTL"

	serviceHost = "SERVICE_HOST"
	servicePort = "SERVICE_PORT"

	shortLinkAddress = "SHORT_LINK_ADDRESS"
)

const (
	_defaultStorageType      = StorageMemory
	_defaultCacheType        = CacheNoop
	_defaultServiceHost      = "localhost"
	_defaultServicePort      = "8080"
	_defaultGenerateAttempts = "10"
	_defaultShortLinkAddress = "http://localhost:8080/link/"
)

const (
	StorageMemory   = "memory"
	StoragePostgres = "postgres"
	CacheNoop       = "noop"
	CacheRedis      = "redis"
)

func InitConfig() *Config {
	envPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("err getting work dir: %v", err.Error()))
	}

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	viper.AutomaticEnv()

	viper.SetDefault(storageType, _defaultStorageType)
	viper.SetDefault(cacheType, _defaultCacheType)
	viper.SetDefault(serviceHost, _defaultServiceHost)
	viper.SetDefault(servicePort, _defaultServicePort)
	viper.SetDefault(generateAttempts, _defaultGenerateAttempts)
	viper.SetDefault(shortLinkAddress, _defaultShortLinkAddress)

	err = viper.ReadInConfig()
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Log.Info(fmt.Sprintf("config file not found: %v", envPath))
		} else {
			panic(fmt.Sprintf("err reading config: %v", err.Error()))
		}
	}

	return &Config{
		PGConfig: &PGConfig{
			PGName:       viper.GetString(PGName),
			PGUser:       viper.GetString(PGUser),
			PGPassword:   viper.GetString(PGPassword),
			PGHost:       viper.GetString(PGHost),
			PGPort:       viper.GetInt(PGPort),
			MaxPool:      viper.GetInt32(PGMaxPool),
			PGTimeout:    viper.GetDuration(PGTimeout),
			ConnAttempts: viper.GetInt(PGConnAttempts),
		},

		RedisConfig: &RedisConfig{
			RedisHost:     viper.GetString(RedisHost),
			RedisPort:     viper.GetInt(RedisPort),
			RedisPassword: viper.GetString(RedisPassword),
		},
		
		StorageType: viper.GetString(storageType),
		CacheType:   viper.GetString(cacheType),

		GenerateAttempts: viper.GetInt(generateAttempts),

		CacheTTL: viper.GetDuration(cacheTTL),

		ServiceHost:      viper.GetString(serviceHost),
		ServicePort:      viper.GetString(servicePort),
		ShortLinkAddress: viper.GetString(shortLinkAddress),
	}
}
