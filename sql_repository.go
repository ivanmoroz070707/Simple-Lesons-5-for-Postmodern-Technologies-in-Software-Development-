package main

import (
	"database/sql"
	"house-api/models" 
	"fmt"
)

type SqlHouseRepository struct {
	db *sql.DB
}

func NewSqlHouseRepository(db *sql.DB) *SqlHouseRepository {
	return &SqlHouseRepository{db: db}
}

func (r *SqlHouseRepository) Create(house *models.House) error {
	query := "INSERT INTO houses (address, price,  rooms, square_meters) VALUES (?, ?, ?, ?)"
	
	result, err := r.db.Exec(query, house.Address, house.Price, house.Rooms, house.SquareMeters)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	house.ID = int(id)
	return nil
}

func (r *SqlHouseRepository) GetByID(id int) (*models.House, error) {
	var h models.House
	query := "SELECT id, address, price,  rooms, square_meters  FROM houses WHERE id = ?"
	
	err := r.db.QueryRow(query, id).Scan(&h.ID, &h.Address, &h.Price, &h.Rooms, &h.SquareMeters)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, fmt.Errorf("помилка пошуку будинку: %w", err)
	}
	
	return &h, nil
}

func (r *SqlHouseRepository) GetAll() ([]models.House, error) {
	query := "SELECT id, address,  price, rooms, square_meters FROM houses"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var houses []models.House
	for rows.Next() {
		var h models.House
		err := rows.Scan(&h.ID, &h.Address, &h.Price, &h.Rooms, &h.SquareMeters)
		if err != nil {
			return nil, err
		}
		houses = append(houses, h)
	}
	return houses, nil
}

func (r *SqlHouseRepository) Delete(id int) error {
	query := "DELETE FROM houses WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("помилка видалення будинку: %w", err)
	}
	return nil
}

func (r *SqlHouseRepository) UpdateFull(house *models.House) error { 
                  query := `UPDATE houses 
	          SET address = ?,  price = ?,  rooms = ?,  square_meters = ? 
	          WHERE id = ?`
	
	_, err := r.db.Exec(query, house.Address, house.Price, house.Rooms, house.SquareMeters, house.ID)
	if err != nil {
		return fmt.Errorf("помилка оновлення будинку: %w", err) }
	return nil

}


func (r *SqlHouseRepository) UpdatePartial(id int, data map[string]interface{}) error { return nil }
