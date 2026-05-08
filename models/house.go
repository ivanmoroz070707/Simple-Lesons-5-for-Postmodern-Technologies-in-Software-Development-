package models

import "fmt"

type House struct {
	ID        int     `json:"id"`
	Address   string  `json:"address"`
	Rooms     int     `json:"rooms"`
	Price     float64 `json:"price"`
}

type HouseRepository interface {
	Create(house *House) error
	GetAll() ([]House, error)
	GetByID(id int) (*House, error)
	UpdateFull(house *House) error
	UpdatePartial(id int, updateData map[string]interface{}) error
	Delete(id int) error
}
func (h *House) Validate() error{
	if strings.TrimSpace(h.Address) == ""{
		return errors.New("addres cannot be empty")
	}
	if h.Price <= 0 {
		return errors.New("price must be a positive number")
	}
	if h.Rooms <= 0 {
		return errors.New("rooms must be a positive number")
	}	
	return nil
}
