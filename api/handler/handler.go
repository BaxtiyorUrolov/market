package handler

import (
	"fmt"
	"market/api/models"
	"market/service"
	"market/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.IStorage
	services service.IServiceManager
}

func New(services service.IServiceManager , store storage.IStorage) Handler {
	return Handler{
		storage: store,
		services: services,
	}
}

func handleResponse(c *gin.Context, msg string, statusCode int, data interface{}) {
	resp := models.Response{}

	switch code := statusCode; {
	case code < 400:
		resp.Description = "success"
	case code < 500:
		resp.Description = "BAD REQUEST"
		fmt.Println("BAD REQUEST:"+msg, " reason: ", data)
	default:
		resp.Description = "INTERNAL SERVER ERROR"
		fmt.Println("INTERVAL SERVER ERROR:"+msg, " reason: ", data)
	}

	resp.StatusCode = statusCode
	resp.Data = data

	c.JSON(resp.StatusCode, resp)
}
