package main

import (
	"database/sql" // Додали цей імпорт, щоб працював sql.Open
	"log"
	"net/http"
	"os"

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

    http.HandleFunc("/houses", handler.CreateHouse)
    
    log.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}

