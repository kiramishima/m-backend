package domain

import "time"

// MarketBond struct
type MarketBond struct {
	ID          int       `json:"id,omitempty" db:"id"`
	UUID        string    `json:"bond_uuid,omitempty" db:"uuid"`
	Name        string    `json:"name,omitempty" db:"name"`
	Price       float32   `json:"price" db:"price"`
	Available   int       `json:"available" db:"available"`
	Currency    int       `json:"currency"  db:"currency"`
	CreatedBy   string    `json:"created_by"  db:"created_by"`
	CreatedByID int       `json:"created_by_id" db:"created_by_id"`
	IsOwner     bool      `json:"is_owner"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
}
