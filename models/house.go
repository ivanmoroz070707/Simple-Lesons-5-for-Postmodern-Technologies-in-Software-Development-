package models

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
