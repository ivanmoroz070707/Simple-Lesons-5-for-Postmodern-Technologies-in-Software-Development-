package main

import (
	"fmt"
	"log"
	"net/http"
	"house-api/models"
)

func main() {
	// 1. Завантаження конфігурації (Viper)
	cfg, err := LoadConfiguration()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	var repo models.HouseRepository
	handler := &HouseHandler{repo: repo} // repo має реалізувати HouseRepository

	// 4. Роутинг (Go 1.22)
	mux := http.NewServeMux()
	

	mux.HandleFunc("POST /houses", handler.CreateHouse)           // 1
	mux.HandleFunc("GET /houses", handler.GetAllHouses)           // 2
	mux.HandleFunc("GET /houses/{id}", handler.GetHouseByID)      // 3
	mux.HandleFunc("PUT /houses/{id}", handler.UpdateHouse)       // 4
	mux.HandleFunc("PATCH /houses/{id}", handler.UpdateHousePartial) // 5
	mux.HandleFunc("DELETE /houses/{id}", handler.DeleteHouse)    // 6

	// 5. Запуск
	fmt.Printf("Server is running on port %s...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}