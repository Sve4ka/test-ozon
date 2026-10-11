package http

import (
	"url_shortener/internal/cache"
	"url_shortener/internal/delivery/http/handler"
	"url_shortener/internal/generate"
	linkService "url_shortener/internal/service/link"
	"url_shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

func initLink(r *gin.Engine, storage storage.Storage, cache cache.Cache, generator generate.Generator, generateAttempt int, shortLinkAddress string) {
	serviceLink := linkService.InitLinkService(storage, cache, generator, generateAttempt)
	handlerLink := handler.InitLinkHandler(serviceLink, shortLinkAddress)

	group := r.Group("/link")
	{
		group.POST("/", handlerLink.Create)
		group.POST("/link/short", handlerLink.GetShort)
		group.GET("/:code", handlerLink.Get)
	}
}

func InitRouter(r *gin.Engine, storage storage.Storage, cache cache.Cache, generator generate.Generator, generateAttempt int, shortLinkAddress string) {
	initLink(r, storage, cache, generator, generateAttempt, shortLinkAddress)
}
