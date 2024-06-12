package service

import (
	"github.com/gorilla/websocket"
	"log"
	"main.go/models"
)

// GetFile - обрабатывает запрос клиента и отправляет данные по WebSocket
func (s *ServiceStruct) GetFile(conn *websocket.Conn, request models.FileRequest) {
	fileChan := make(chan []byte)
	errChan := make(chan error)

	go s.Decompress(request.Login, request.Password, request.ID, fileChan, errChan)

	select {
	case file := <-fileChan:
		err := conn.WriteMessage(websocket.BinaryMessage, file)
		if err != nil {
			log.Println("Error sending file over WebSocket:", err)
		}
	case err := <-errChan:
		log.Println("Error:", err)
	}
}

func (s *ServiceStruct) PutFile(record models.Record) (string, error) {
	file, err := s.СompressFile(record.File)
	if err != nil {
		log.Println("Error of compress file:", err)
		return "", err
	}
	record.File = file
	s.storage.WriteRecord(record)
	return "", nil
}

func (s *ServiceStruct) DeleteFile(id string, login string, password string) (string, error) {
	return "", nil
}
