package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/serjbibox/SF36a.3/pkg/storage"
)

// Обработчик HTTP запросов сервера GoNews
type Handler struct {
	storage *storage.Storage
}

//Конструктор объекта Handler
func New(storage *storage.Storage) (*Handler, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}
	return &Handler{storage: storage}, nil
}

//Инициализация маршрутизатора запросов.
//Регистрация обработчиков запросов
func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./cmd/news/webapp", true)))
	r.GET("/news/:n", h.getNews)
	return r
}

// Получение публикаций по заданному количеству
func (h *Handler) getNews(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	posts, err := h.storage.Post.GetByQuantity(n)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, posts)
}
