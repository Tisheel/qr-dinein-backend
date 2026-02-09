package model

import "time"

type Category struct {
	ID           int       `json:"id"`
	RestaurantID int       `json:"restaurantId"`
	Name         string    `json:"name"`
	Order        int       `json:"order"`
	Image        string    `json:"image"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
