package api

import (
	"github.com/gin-gonic/gin"
	"main.go/models"
	"net/http"
)

// структура для хранения информации о файле
type FileInfo struct {
	FileName string
	FilePath string
}

// структура для передачи информации о файлах и учетных данных
type FileTransferRequest struct {
	Files    []FileInfo
	Username string
	Password string
}

func (h *MyHandler) UploadFile(c *gin.Context) {
	var request models.FileRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание экземпляра Record из запроса
	record := models.Record{
		ID:       request.ID,
		Login:    request.Login,
		Password: request.Password,
	}

	// Вызов метода PutFile из ServiceStruct
	_, err := h.MyService.PutFile(record)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File(s) uploaded successfully"})
}

func (h *MyHandler) DownloadFile(c *gin.Context) {
	// Обработка скачивания файла
	// Пример получения данных о файлах и учетных данных
	var request FileTransferRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// В request.Files будут содержаться информация о файлах, которые необходимо скачать
	// В request.Username и request.Password будут содержаться учетные данные

	// Обработка скачивания файла...

	c.JSON(http.StatusOK, gin.H{"message": "File(s) downloaded successfully"})
}
