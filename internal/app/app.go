package app

import (
	"backend/internal/cache"
	"backend/internal/cache/noop"
	"backend/internal/cache/redis"
	"backend/internal/config"
	"backend/internal/delivery/http"
	"backend/internal/delivery/middleware"
	"backend/internal/docs"
	"backend/internal/generate"
	"backend/internal/log"
	"backend/internal/storage"
	"backend/internal/storage/memory"
	"backend/internal/storage/postgres"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start() {
	cfg := config.InitConfig()

	log.Log.Info("Config Initialized")

	var db storage.Storage

	switch cfg.StorageType {
	case config.StorageMemory:
		inMemory := memory.NewStorage()
		db = memory.InitMemoryRepo(inMemory)
	case config.StoragePostgres:
		pg := postgres.MustInitPg(cfg.PGConfig)
		db = postgres.InitPostgresRepo(pg)
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

	http.InitRouter(g, db, cacheData, generator, cfg.GenerateAttempts, cfg.ShortLinkAddress)

	log.Log.Info("Start Server")

	err := g.Run(fmt.Sprintf("%v:%v", cfg.ServiceHost, cfg.ServicePort))
	if err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
