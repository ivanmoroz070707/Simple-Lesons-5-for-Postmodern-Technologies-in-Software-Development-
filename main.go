package main

import (
	"database/sql" // Додали цей імпорт, щоб працював sql.Open
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql" // Драйвер MySQL
)

type House struct {
    ID                          string  `json:"id"`
    Address              string  `json:"address"`
    Price                     int     `json:"price"`
    Floors                  int     `json:"floors"`       
    SquareMeters  float64 `json:"square_meters"` 
}

func main() {
	// 1. Завантаження конфігурації
	cfg, err := LoadConfiguration()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// 2. Підключення до MySQL
	db, err := sql.Open("mysql", cfg.DBURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Перевірка зв'язку
	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	// 3. Запуск міграцій (зверни увагу на велику літеру R)
	RunMigrations(db)

	// 4. Роутер
	r := chi.NewRouter()
	
	// Тут Ваня  додай свої обробники (handlers)

	fmt.Printf("Server starting on port %s...\n", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}
