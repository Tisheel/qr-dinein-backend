package model

import "time"

type Product struct {
	ID           int       `json:"id"`
	RestaurantID int       `json:"restaurantId"`
	CategoryID   int       `json:"categoryId"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	Image        string    `json:"image"`
	Veg          bool      `json:"veg"`
	Available    bool      `json:"available"`
	PrepTime     int       `json:"prepTime"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
