package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main.go/config"
	"main.go/internal/api"
	"main.go/internal/service"
	"main.go/internal/storage"
	"main.go/server"
	"os"
	"time"
)

var Logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}

	Logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func MigrateDatabase(dbConfig config.DatabaseConfig) (*sql.DB, error) {
	// Устанавливаем соединение с базой данных SQLite
	db, err := sql.Open(dbConfig.Driver, dbConfig.Source)
	if err != nil {
		return nil, err
	}

	// Проверяем, существует ли файл базы данных
	if _, err := os.Stat(dbConfig.Source); err == nil {
		// Если файл существует, возвращаем указатель на sql.DB
		return db, sql.ErrConnDone
	}

	// Создаем файл базы данных
	file, err := os.Create(dbConfig.Source)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Создаем таблицу users
	_, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            password TEXT NOT NULL,
            file BLOB,
            hash TEXT
        )`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	InitLogger()

	Config := config.Conf

	db, err := MigrateDatabase(Config.Database)
	if err != nil {
		if err == sql.ErrConnDone {
			log.Println("База данных уже существует")
		} else {
			log.Fatal("Ошибка создания базы данных:", err)
		}
	} else {
		log.Println("База данных успешно создана и подключена")
	}
	defer func() {
		log.Println("Закрываю соединение с бд")
		err := db.Close()
		if err != nil {
			log.Println("Ошибка закрытия соединения")
		}
	}()

	Mystorage := storage.NewStorage(db)
	MyService := service.NewService(Mystorage)
	MyHandler := api.NewMyHandler(&MyService)

	server := new(server.Server)

	Logger.Println("Starting server on port", config.Conf.Server.HTTP.Port)
	if Config.Server.HTTP.ScheduledShutdown == 0 {

	} else {
		ctx, cancel := context.WithTimeout(context.Background(), Config.Server.HTTP.ScheduledShutdown)
		defer cancel() //на случай принудительного завершения

		go func(ctx context.Context) {
			if err := server.RunServer(MyHandler.InitRouts()); err != nil {
				log.Fatal("Server start error: ", err)
			}
		}(ctx)

		time.AfterFunc(Config.Server.HTTP.ScheduledShutdown, cancel)
	}

	fmt.Scanln()
}
