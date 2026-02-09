package service

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/store"

	"gofr.dev/pkg/gofr"
)

type RatingService struct {
	store      *store.Rating
	orderStore *store.Order
}

func NewRating(s *store.Rating, orderStore *store.Order) *RatingService {
	return &RatingService{store: s, orderStore: orderStore}
}

func (svc *RatingService) GetByOrderID(ctx *gofr.Context, restaurantID, orderID int) (*model.Rating, error) {
	return svc.store.GetByOrderID(ctx, restaurantID, orderID)
}

func (svc *RatingService) GetAllByRestaurant(ctx *gofr.Context, restaurantID int) ([]model.Rating, error) {
	return svc.store.GetAllByRestaurant(ctx, restaurantID)
}

func (svc *RatingService) Create(ctx *gofr.Context, restaurantID, orderID int, r *model.Rating) (*model.Rating, error) {
	// Validate rating value
	if r.Rating < 1 || r.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	// Verify order exists and belongs to the restaurant
	order, err := svc.orderStore.GetByID(ctx, restaurantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Only completed orders can be rated
	if order.Status != "completed" {
		return nil, fmt.Errorf("only completed orders can be rated")
	}

	// Check if already rated
	existing, _ := svc.store.GetByOrderID(ctx, restaurantID, orderID)
	if existing != nil {
		return nil, fmt.Errorf("order has already been rated")
	}

	r.OrderID = orderID
	r.RestaurantID = restaurantID

	return svc.store.Create(ctx, r)
}
