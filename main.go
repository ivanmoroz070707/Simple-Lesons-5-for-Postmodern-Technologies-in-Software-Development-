package main

import (
	"database/sql" // Додали цей імпорт, щоб працював sql.Open
	"log"
	"net/http"
	"os"
	"context"
	"syscall"
	"time"
	"os/signal"

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
    // 1. Підключення до бази 
    db, err := sql.Open("mysql", os.Getenv("DB_URL"))
    if err != nil {  log.Fatal(err) }

    houseRepo := NewSqlHouseRepository(db)
    handler := &HouseHandler{ repo: houseRepo}
    
    r := chi.NewRouter()
    
    
    
    r.Post("/houses", handler.CreateHouse)
    r.Get("/houses", handler.GetHouses)
    
    r.Get("/houses/{id}", handler.GetHouseByID)
    r.Put("/houses/{id}", handler.UpdateHouse)    
    r.Delete("/houses/{id}", handler.DeleteHouse) 

    http.HandleFunc("/houses", handler.CreateHouse)
    
    log.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
    
    srv := &http.Server{
		Addr:    ":8080",
		Handler: r, // твій роутер chi
	}

	// 2. Запускаємо сервер у фоновому потоці (goroutine)
	go func() {
		log.Println("Сервер запущено на порту 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Помилка сервера: %s\n", err)
		}
	}()

	// 3. Чекаємо на сигнал зупинки від системи (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Код зупиняється тут і чекає...

	log.Println("Отримано сигнал зупинки. Вимикаємо сервер...")

	// 4. Даємо серверу 5 секунд, щоб завершити поточні запити і відпустити порт
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Примусове вимкнення:", err)
	}

	log.Println("Сервер успішно зупинено. Порт вільний!")
}

