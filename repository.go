package main

import "house-api/models"

type HouseRepository interface {
	Create(house *models.House) error
	GetAll() ([]models.House,  error)
	GetByID(id int) (*models.House, error)
	UpdateFull(house *models.House) error
	UpdatePartial(id int, updateData map[string]interface{}) error
	Delete(id int) error
}
