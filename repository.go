package main

type HouseRepository interface {
	Create(house *House) error
	GetAll() ([]House, error)
	GetByID(id string) (*House, error)
	UpdateFull(house *House) error
	UpdatePartial(id string, updateData map[string]interface{}) error
	Delete(id string) error
}
