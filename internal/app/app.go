package app

import (
	"fmt"
	"url_shortener/internal/cache"
	"url_shortener/internal/cache/noop"
	"url_shortener/internal/cache/redis"
	"url_shortener/internal/config"
	"url_shortener/internal/delivery/http"
	"url_shortener/internal/delivery/middleware"
	"url_shortener/internal/docs"
	"url_shortener/internal/generate"
	"url_shortener/internal/log"
	"url_shortener/internal/storage"
	"url_shortener/internal/storage/memory"
	"url_shortener/internal/storage/postgres"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start() {
	cfg := config.InitConfig()

	log.Log.Info("Config Initialized")

	var store storage.Storage

	switch cfg.StorageType {
	case config.StorageMemory:
		inMemory := memory.NewStorage()
		store = memory.InitMemoryRepo(inMemory)
	case config.StoragePostgres:
		pg := postgres.MustInitPg(cfg.PGConfig)
		store = postgres.InitPostgresRepo(pg)
		defer pg.Close()
	default:
		panic("Unsupported storage type")
	}

	var cacheData cache.Cache

	switch cfg.CacheType {
	case config.CacheNoop:
		cacheData = noop.InitNoopCache()
	case config.CacheRedis:
		redisClient, err := redis.InitRedis(cfg.RedisConfig)
		if err != nil {
			log.Log.Error(err)
			cacheData = noop.InitNoopCache()
		} else {
			cacheData = redis.InitRedisCache(redisClient, cfg.CacheTTL)
		}
	default:
		panic("unsupported cache type")
	}

	log.Log.Info("Storage Initialized")

	generator := generate.NewRandomGenerator()
	middlewareStruct := middleware.InitMiddleware()

	g := gin.New()
	g.Use(middlewareStruct.CORSMiddleware())

	docs.SwaggerInfo.BasePath = "/"
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	http.InitRouter(g, store, cacheData, generator, cfg.GenerateAttempts, cfg.ShortLinkAddress)

	log.Log.Info("Start Server")

	err := g.Run(fmt.Sprintf("%v:%v", cfg.ServiceHost, cfg.ServicePort))
	if err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
