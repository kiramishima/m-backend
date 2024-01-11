package domain

import "time"

// Bond struct
type Bond struct {
	ID          int       `json:"id,omitempty" db:"id"`
	UUID        string    `json:"bond_id,omitempty" db:"uuid"`
	Name        string    `json:"name,omitempty" db:"name"`
	Price       float32   `json:"price" db:"price"`
	Number      int       `json:"num" db:"number"`
	Currency    int       `json:"currency"  db:"currency"`
	CreatedBy   string    `json:"created_by"  db:"created_by"`
	CreatedByID int       `json:"created_by_id" db:"created_by_id"`
	OnSale      bool      `json:"on_sale" db:"on_sale"`
	IsOwner     bool      `json:"is_owner"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
}
