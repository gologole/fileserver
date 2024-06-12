package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"main.go/models"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *MyHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	var request models.FileRequest
	err = conn.ReadJSON(&request)
	if err != nil {
		log.Println("Error reading message from WebSocket:", err)
		return
	}

	h.MyService.GetFile(conn, request)
}
