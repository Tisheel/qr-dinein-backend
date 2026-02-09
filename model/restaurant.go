package model

import "time"

type Restaurant struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Logo      string    `json:"logo"`
	Currency  string    `json:"currency"`
	TaxRate   float64   `json:"taxRate"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
