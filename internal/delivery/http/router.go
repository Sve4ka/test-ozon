package http

import (
	"backend/internal/cache"
	"backend/internal/delivery/http/handler"
	"backend/internal/generate"
	linkService "backend/internal/service/link"
	"backend/internal/storage"

	"github.com/gin-gonic/gin"
)

func initLink(r *gin.Engine, db storage.Storage, cache cache.Cache, generator generate.Generator, generateAttempt int, shortLinkAddress string) {
	serviceLink := linkService.InitLinkService(db, cache, generator, generateAttempt)
	handlerLink := handler.InitLinkHandler(serviceLink, shortLinkAddress)

	group := r.Group("/link")
	{
		group.POST("/", handlerLink.Create)
		group.POST("/link/short", handlerLink.GetShort)
		group.GET("/:code", handlerLink.Get)
	}
}

func InitRouter(r *gin.Engine, db storage.Storage, cache cache.Cache, generator generate.Generator, generateAttempt int, shortLinkAddress string) {
	initLink(r, db, cache, generator, generateAttempt, shortLinkAddress)
}
