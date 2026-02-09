package model

type Setting struct {
	ID           int    `json:"id"`
	RestaurantID int    `json:"restaurantId"`
	Key          string `json:"key"`
	Value        string `json:"value"`
}
