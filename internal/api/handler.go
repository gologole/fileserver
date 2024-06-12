package api

import (
	"github.com/gin-gonic/gin"
	"main.go/internal/service"
)

type MyHandler struct {
	MyService service.Service
}

func NewMyHandler(service service.ServiceStruct) *MyHandler {
	return &MyHandler{
		MyService: &service,
	}
}

func (h *MyHandler) InitRouts() *gin.Engine {
	router := gin.Default()

	// Раздача HTML файла из корня проекта
	router.GET("/site", func(c *gin.Context) {
		c.File("index.html")
	})

	router.Group("/httpp")
	{
		router.POST("/uploadfile", h.UploadFile)
		router.POST("/downloadfile", h.DownloadFile)
	}

	router.GET("/ws", h.HandleWebSocket)

	return router
}
