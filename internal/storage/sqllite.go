package storage

import (
	"database/sql"
	"fmt"
	"main.go/models"
	"strings"
	"sync"
)

type Storage interface {
	Open() error
	GetFilePartByID(login, password, id string, partChan chan<- []byte, errChan chan<- error)
	WriteRecord(record models.Record) error
	GetHashAndMetadataByID(id string) (string, []string, error)
	DeleteRecord(id string) error
	Close() error
}

type storage struct {
	db *sql.DB
	mu *sync.Mutex
}

func NewStorage(db *sql.DB) *storage {
	return &storage{
		db: db,
		mu: &sync.Mutex{},
	}
}

func (s *storage) Open() error {
	return nil
}

func (s *storage) GetFilePartByID(login, password, id string, partChan chan<- []byte, errChan chan<- error) {
	// Блокировка доступа к данным на время выполнения запроса
	s.mu.Lock()
	defer s.mu.Unlock()

	// SQL-запрос для извлечения файла по ID и проверки логина и пароля
	query := "SELECT File, Login, Password FROM records WHERE ID=?"
	row := s.db.QueryRow(query, id)

	var file []byte
	var retrievedLogin, retrievedPassword string
	err := row.Scan(&file, &retrievedLogin, &retrievedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			errChan <- fmt.Errorf("record not found")
		} else {
			errChan <- fmt.Errorf("error retrieving record: %v", err)
		}
		close(partChan)
		return
	}

	// Проверка логина и пароля
	if login != retrievedLogin || password != retrievedPassword {
		errChan <- fmt.Errorf("invalid login/password")
		close(partChan)
		return
	}

	// Отправка частей файла через канал
	partSize := 1024 // Размер частей файла
	for i := 0; i < len(file); i += partSize {
		end := i + partSize
		if end > len(file) {
			end = len(file)
		}
		partChan <- file[i:end]
	}
	close(partChan)
}

func (s *storage) WriteRecord(record models.Record) error {
	// Блокировка доступа к данным на время выполнения запроса
	s.mu.Lock()
	defer s.mu.Unlock()

	// Подготовка SQL-запроса для вставки записи в базу данных
	query := "INSERT INTO records (ID, Login, Password, Metadata, File, Hash) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := s.db.Exec(query, record.ID, record.Login, record.Password, record.Metadata, record.File, record.Hash)
	if err != nil {
		return fmt.Errorf("error writing record to database: %v", err)
	}

	return nil
}

func (s *storage) GetHashAndMetadataByID(id string) (string, []string, error) {
	// Блокировка доступа к данным на время выполнения запроса
	s.mu.Lock()
	defer s.mu.Unlock()

	// Подготовка SQL-запроса для получения хеша и метаданных по ID
	query := "SELECT Hash, Metadata FROM records WHERE ID=?"
	row := s.db.QueryRow(query, id)

	var hash string
	var metadataString string
	err := row.Scan(&hash, &metadataString)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, fmt.Errorf("record not found")
		}
		return "", nil, fmt.Errorf("error retrieving hash and metadata: %v", err)
	}

	// Разбиваем строку метаданных на отдельные элементы
	metadata := strings.Split(metadataString, ",")

	return hash, metadata, nil
}

func (s *storage) DeleteRecord(id string) error {
	// Блокировка доступа к данным на время выполнения запроса
	s.mu.Lock()
	defer s.mu.Unlock()

	// Подготовка SQL-запроса для удаления записи по ID
	query := "DELETE FROM records WHERE ID=?"
	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting record: %v", err)
	}

	return nil
}

func (s *storage) Close() error {
	// Закрытие соединения с базой данных
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %v", err)
	}
	return nil
}
