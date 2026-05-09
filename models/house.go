package models

import "errors"

type House struct {
	ID           int     `json:"id"`
	Address      string  `json:"address"`
	Rooms        int     `json:"rooms"`
	Price        float64 `json:"price"`
	SquareMeters float64 `json:"square_meters"` 
}


func (h *House) Validate() error {
	if h.Address == "" {
		return errors.New("адреса не може бути порожньою")
	}
	if h.Price <= 0 {
		return errors.New("ціна має бути більшою за 0")
	}
	return nil
}

type HouseRepository interface {
	Create(house *House) error
	GetAll() ([]House, error)
	GetByID(id int) (*House, error)
	UpdateFull(house *House) error
	UpdatePartial(id int, updateData map[string]interface{}) error
	Delete(id int) error
}
