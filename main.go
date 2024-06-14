package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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

func MigrateDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS records (
        id TEXT PRIMARY KEY,
        login TEXT,
        password TEXT,
        metadata TEXT,
        file BLOB,
        hash TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	InitLogger()

	Config := config.Conf

	db, err := MigrateDatabase(Config.Database.Source)
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

	fmt.Println("инициализация хранилища")
	Mystorage := storage.NewStorage(db)
	fmt.Println("инициализация сервиса")
	MyService := service.NewService(Mystorage)
	fmt.Println("инициализация хендлеров")
	MyHandler := api.NewMyHandler(MyService)

	fmt.Println("инициализация сервера")
	server := new(server.Server)

	Logger.Println("Starting server on port", config.Conf.Server.HTTP.Port)

	if Config.Server.HTTP.ScheduledShutdown == 0 {
		fmt.Println("Запуск сервера")
		if err := server.RunServer(MyHandler.InitRouts()); err != nil {
			log.Fatal("Server start error: ", err)
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), Config.Server.HTTP.ScheduledShutdown)
		defer cancel() //на случай принудительного завершения

		go func(ctx context.Context) {
			fmt.Println("Запуск сервера c запланированным окончанием работы")
			if err := server.RunServer(MyHandler.InitRouts()); err != nil {
				log.Fatal("Server start error: ", err)
			}
		}(ctx)

		time.AfterFunc(Config.Server.HTTP.ScheduledShutdown, cancel)
	}

	fmt.Scanln()
}
