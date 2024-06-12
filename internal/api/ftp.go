package api

import (
	ftpserver "github.com/fclairamb/ftpserverlib"
	"log"
	"os"
	"path/filepath"
)

// User - структура для хранения информации о пользователе
type User struct {
	Username string
	Password string
	BaseDir  string
}

// FTPDriver - реализация интерфейса ftpserverlib.Driver
type FTPDriver struct {
	users map[string]User
}

// NewFTPDriver создает новый FTPDriver
func NewFTPDriver(users map[string]User) *FTPDriver {
	return &FTPDriver{users: users}
}

// WelcomeUser - приветственное сообщение
func (driver *FTPDriver) WelcomeUser(cc ftpserver.ClientContext) (string, error) {
	return "Welcome to the FTP server", nil
}

// AuthUser - аутентификация пользователя
func (driver *FTPDriver) AuthUser(cc ftpserver.ClientContext, user, pass string) (ftpserver.ClientHandlingDriver, error) {
	if u, ok := driver.users[user]; ok && u.Password == pass {
		return &ClientHandler{
			BaseDir: u.BaseDir,
		}, nil
	}
	return nil, ftpserver.ErrInvalidLogin
}

// ClientHandler - обработчик FTP команд
type ClientHandler struct {
	BaseDir string
}

// ChangeDirectory - смена директории
func (handler *ClientHandler) ChangeDirectory(cc ftpserver.ClientContext, directory string) error {
	return nil
}

// MakeDirectory - создание директории
func (handler *ClientHandler) MakeDirectory(cc ftpserver.ClientContext, directory string) error {
	path := handler.fullPath(directory)
	return os.Mkdir(path, 0755)
}

// ListFiles - список файлов в директории
func (handler *ClientHandler) ListFiles(cc ftpserver.ClientContext) ([]os.FileInfo, error) {
	path := handler.fullPath(cc.Path())
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []os.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, info)
	}

	return files, nil
}

// OpenFile - открытие файла
func (handler *ClientHandler) OpenFile(cc ftpserver.ClientContext, path string, flag int) (ftpserver.FileStream, error) {
	fullPath := handler.fullPath(path)
	f, err := os.OpenFile(fullPath, flag, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// DeleteFile - удаление файла
func (handler *ClientHandler) DeleteFile(cc ftpserver.ClientContext, path string) error {
	fullPath := handler.fullPath(path)
	return os.Remove(fullPath)
}

// RenameFile - переименование файла
func (handler *ClientHandler) RenameFile(cc ftpserver.ClientContext, from, to string) error {
	fromPath := handler.fullPath(from)
	toPath := handler.fullPath(to)
	return os.Rename(fromPath, toPath)
}

// GetFileInfo - информация о файле
func (handler *ClientHandler) GetFileInfo(cc ftpserver.ClientContext, path string) (os.FileInfo, error) {
	fullPath := handler.fullPath(path)
	return os.Stat(fullPath)
}

// fullPath - возвращает полный путь относительно базовой директории пользователя
func (handler *ClientHandler) fullPath(path string) string {
	cleanPath := filepath.Clean(path)
	return filepath.Join(handler.BaseDir, cleanPath)
}

// StartFTPServer запускает FTP сервер
func StartFTPServer(port string, users map[string]User) {
	// Создаем базовую директорию
	for _, user := range users {
		os.MkdirAll(user.BaseDir, 0755)
	}

	// Создаем FTP драйвер
	driver := NewFTPDriver(users)

	// Конфигурируем и запускаем сервер
	server := ftpserver.NewFtpServer(driver)
	server.Logger = log.New(os.Stdout, "ftpserver: ", log.LstdFlags)

	log.Printf("Starting FTP server on %s\n", port)
	if err := server.ListenAndServe(port); err != nil {
		log.Fatal(err)
	}
}
