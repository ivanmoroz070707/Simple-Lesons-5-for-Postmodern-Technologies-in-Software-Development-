package main

import("context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool")
type House struct{
	ID      string `json:"id"`
	Address string `json:"address"`
	Price   int    `json:"price"`
}
func main(){
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	dbpool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	r:=chi.NewRouter()


	fmt.Printf("Server starting on port %s...\n", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}