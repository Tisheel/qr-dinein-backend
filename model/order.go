package model

import "time"

type OrderItem struct {
	ProductID int     `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Veg       bool    `json:"veg"`
}

type Order struct {
	ID                  int        `json:"id"`
	RestaurantID        int        `json:"restaurantId"`
	TableNumber         *string    `json:"tableNumber"`
	CustomerMobile      string     `json:"customerPhone"`
	CustomerName        string     `json:"customerName"`
	Items               []OrderItem `json:"items"`
	Status              string     `json:"status"`
	SpecialInstructions string     `json:"specialInstructions"`
	Total               float64    `json:"total"`
	AssignedChefID      *int       `json:"assignedChefId"`
	EstimatedReadyAt    *time.Time `json:"estimatedReadyAt"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`

	// SessionToken is used only for customer order creation (not stored in DB)
	SessionToken string `json:"sessionToken,omitempty"`
}
