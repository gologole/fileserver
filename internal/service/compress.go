package service

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

func (s *ServiceStruct) СompressFile(data []byte) ([]byte, error) {
	// Создаем буфер для записи сжатых данных
	var compressedBuffer bytes.Buffer

	// Создаем новый zlib.Writer с буфером
	zw := zlib.NewWriter(&compressedBuffer)

	// Записываем данные в zlib.Writer
	_, err := zw.Write(data)
	if err != nil {
		return nil, err
	}

	// Закрываем zlib.Writer
	err = zw.Close()
	if err != nil {
		return nil, err
	}

	// Возвращаем сжатые данные
	return compressedBuffer.Bytes(), nil
}

// Функция декомпрессии части файла
func decompressPart(part []byte, decompressedChan chan<- []byte, errChan chan<- error) {
	// Создаем буфер для записи декомпрессированных данных
	var decompressedBuffer bytes.Buffer

	// Создаем новый zlib.Reader с частью данных
	zr, err := zlib.NewReader(bytes.NewReader(part))
	if err != nil {
		errChan <- fmt.Errorf("error creating zlib reader: %v", err)
		return
	}
	defer zr.Close()

	// Копируем данные из zlib.Reader в буфер
	_, err = io.Copy(&decompressedBuffer, zr)
	if err != nil {
		errChan <- fmt.Errorf("error decompressing part: %v", err)
		return
	}

	// Отправляем декомпрессированные данные через канал
	decompressedChan <- decompressedBuffer.Bytes()
}

// Функция декомпрессии с параллельным чтением из базы данных
func (s *ServiceStruct) Decompress(login, password, id string, fileChan chan<- []byte, errChan chan<- error) {
	// Создаем каналы для передачи частей данных и ошибок
	partChan := make(chan []byte)
	decompressedChan := make(chan []byte)
	dbErrChan := make(chan error)

	// Запускаем функцию чтения из базы данных в отдельной горутине
	go func() {
		s.storage.GetFilePartByID(login, password, id, partChan, dbErrChan)
		close(partChan) // Закрываем канал после завершения работы
	}()

	// Запускаем горутину для обработки ошибок базы данных
	go func() {
		if err := <-dbErrChan; err != nil {
			errChan <- err
			close(errChan)
			return
		}
	}()

	// Запускаем горутину для декомпрессии данных
	go func() {
		for part := range partChan {
			decompressPart(part, decompressedChan, errChan)
		}
		close(decompressedChan) // Закрываем канал после завершения работы
	}()

	// Объединяем декомпрессированные части и отправляем через fileChan
	var decompressedFile bytes.Buffer
	for decompressedPart := range decompressedChan {
		decompressedFile.Write(decompressedPart)
	}
	fileChan <- decompressedFile.Bytes()
	close(fileChan)
}
