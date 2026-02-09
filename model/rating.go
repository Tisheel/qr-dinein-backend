package model

import "time"

type Rating struct {
	ID           int       `json:"id"`
	OrderID      int       `json:"orderId"`
	RestaurantID int       `json:"restaurantId"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"createdAt"`
}
