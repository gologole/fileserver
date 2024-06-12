package models

import "github.com/gorilla/websocket"

// FileRequest - структура, ожидаемая от клиента
type FileRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Record struct {
	ID       string
	Login    string
	Password string
	Metadata []string
	File     []byte
	Hash     string
}
