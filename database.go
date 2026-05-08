package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" 
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("Помилка створення драйвера міграцій: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", 
		"mysql", driver)
	if err != nil {
		log.Fatalf("Помилка ініціалізації міграцій: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Помилка під час виконання міграцій: %v", err)
	}

	log.Println("Міграції успішно перевірені/застосовані!")
}

