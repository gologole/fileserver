package service

import (
	"github.com/gorilla/websocket"
	"io"
	"main.go/internal/storage"
	"main.go/models"
)

// компрессию и декомпрессию сделать конвеером горутин
type Service interface {
	Compress
	Cache
	GetFile(conn *websocket.Conn, request models.FileRequest)
	PutFile(record models.Record) (string, error)
	DeleteFile(id string, login string, password string) (string, error)
}

type Compress interface {
	Decompress(login, password, id string, fileChan chan<- []byte, errChan chan<- error)
	СompressFile(data []byte) ([]byte, error)
}

type Cache interface {
	getcache() (io.ReadCloser, error)
	putcache() (io.ReadCloser, error)
	deletecache() (io.ReadCloser, error)
}

type ServiceStruct struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) Service {
	return &ServiceStruct{storage}
}
