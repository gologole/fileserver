package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"main.go/internal/api"
	"main.go/internal/service"
	"main.go/models"
	"net/http/httptest"
	"testing"
)

func TestWebSocketEndpoint(t *testing.T) {
	// Создаем Gin router
	router := gin.Default()

	// Инициализируем ваш хэндлер и сервис
	myHandler := &api.MyHandler{
		MyService: &service.ServiceStruct{},
	}

	// Регистрируем WebSocket эндпоинт
	router.GET("/ws", myHandler.HandleWebSocket)

	// Создаем тестовый сервер
	server := httptest.NewServer(router)
	defer server.Close()

	// Получаем URL для подключения к WebSocket
	u := "ws" + server.URL[len("http"):]

	// Устанавливаем WebSocket соединение
	ws, _, err := websocket.DefaultDialer.Dial(u+"/ws", nil)
	require.NoError(t, err)
	defer ws.Close()

	// Отправляем тестовый запрос
	request := models.FileRequest{
		Login:    "test_login",
		Password: "test_password",
		ID:       "test_id",
	}
	err = ws.WriteJSON(request)
	require.NoError(t, err)

	// Читаем ответ
	_, message, err := ws.ReadMessage()
	require.NoError(t, err)

	// Проверяем, что получили ожидаемый ответ
	expectedFileContent := []byte("expected file content")
	assert.Equal(t, expectedFileContent, message)
}
