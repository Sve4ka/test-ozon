package handler

import (
	"backend/internal/cerr"
	"backend/internal/generate"
	"backend/internal/models"
	"backend/internal/service"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Link struct {
	service          service.Link
	shortLinkAddress string
}

func InitLinkHandler(service service.Link, shortLinkAddress string) *Link {
	return &Link{
		service:          service,
		shortLinkAddress: shortLinkAddress,
	}
}

// @Summary Создает короткую ссылку по оригинальной
// @Tags Link
// @Accept	json
// @Produce	json
// @Param link body models.OriginalLinkRequest true "OriginalLink"
// @Success 201 {object} models.ShortCodeResponse "Ссылка создана"
// @Failure 400 {object} models.Error "Неверный запрос"
// @Failure 500 {object} models.Error "Внутренняя ошибка сервера"
// @Router /link/ [post]
func (s Link) Create(c *gin.Context) {
	var request models.OriginalLinkRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Неверный запрос"})

		return
	}

	parse, err := url.ParseRequestURI(string(request.OriginalLink))
	if err != nil || (parse.Scheme != "http" && parse.Scheme != "https") || parse.Host == "" {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Неверный запрос"})

		return
	}

	ctx := c.Request.Context()
	code, err := s.service.Create(ctx, request.OriginalLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Error: "Внутренняя ошибка сервера"})

		return
	}
	c.JSON(http.StatusCreated, models.ShortCodeResponse{
		ShortCode: *code,
		ShortLink: s.shortLinkAddress + string(*code),
	})
}

// @Summary Возвращает оригинальную ссылку по короткой ссылке
// @Tags Link
// @Accept	json
// @Produce	json
// @Param link body models.ShortCodeRequest true "OriginalLink"
// @Success 200 {object} models.OriginalLinkResponse "Оригинальная ссылка"
// @Failure 400 {object} models.Error "Неверный запрос"
// @Failure 404 {object} models.Error "Ссылка не найдена"
// @Failure 500 {object} models.Error "Внутренняя ошибка сервера"
// @Router /link/short [post]
func (s Link) GetShort(c *gin.Context) {
	var request models.ShortCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Неверный запрос"})

		return
	}

	parse, err := url.ParseRequestURI(string(request.ShortLink))
	if err != nil || (parse.Scheme != "http" && parse.Scheme != "https") || parse.Host == "" {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Неверный запрос"})

		return
	}
	code, found := strings.CutPrefix(request.ShortLink, s.shortLinkAddress)
	if !generate.IsValidCode(code) || !found {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Неверный запрос"})

		return
	}

	ctx := c.Request.Context()
	originalLink, err := s.service.Get(ctx, models.ShortCode(code))
	if err != nil {
		if errors.Is(err, cerr.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.Error{Error: "Ссылка не найдена"})

			return
		}
		c.JSON(http.StatusInternalServerError, models.Error{Error: "Внутренняя ошибка сервера"})

		return
	}
	c.JSON(http.StatusOK, models.OriginalLinkResponse{OriginalLink: *originalLink})
}

// @Summary Возвращает оригинальную ссылку по короткому коду
// @Tags Link
// @Accept	json
// @Produce	json
// @Param code path string true "Содержание короткой ссылки"
// @Success 200 {object} models.OriginalLinkResponse "Оригинальная ссылка"
// @Failure 400 {object} models.Error "Неверный запрос"
// @Failure 404 {object} models.Error "Ссылка не найдена"
// @Failure 500 {object} models.Error "Внутренняя ошибка сервера"
// @Router /link/{code} [get]
func (s Link) Get(c *gin.Context) {
	code := c.Param("code")

	if !generate.IsValidCode(code) {
		c.JSON(http.StatusBadRequest, models.Error{Error: "Неверный запрос"})

		return
	}

	ctx := c.Request.Context()
	originalLink, err := s.service.Get(ctx, models.ShortCode(code))
	if err != nil {
		if errors.Is(err, cerr.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.Error{Error: "Ссылка не найдена"})

			return
		}
		c.JSON(http.StatusInternalServerError, models.Error{Error: "Внутренняя ошибка сервера"})

		return
	}
	c.JSON(http.StatusOK, models.OriginalLinkResponse{OriginalLink: *originalLink})
}
