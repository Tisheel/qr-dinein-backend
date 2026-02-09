package model

import "time"

type Staff struct {
	ID           int       `json:"id"`
	RestaurantID int       `json:"restaurantId"`
	Username     string    `json:"username"`
	Pin          string    `json:"pin,omitempty"`
	Role         string    `json:"role"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
