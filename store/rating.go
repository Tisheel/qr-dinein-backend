package store

import (
	"qr-dinein-backend/model"
	"time"

	"gofr.dev/pkg/gofr"
)

type Rating struct{}

func NewRating() *Rating {
	return &Rating{}
}

func (s *Rating) GetByOrderID(ctx *gofr.Context, restaurantID, orderID int) (*model.Rating, error) {
	var r model.Rating
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, order_id, restaurant_id, rating, comment, created_at FROM order_ratings WHERE order_id = ? AND restaurant_id = ?",
		orderID, restaurantID).
		Scan(&r.ID, &r.OrderID, &r.RestaurantID, &r.Rating, &r.Comment, &r.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *Rating) GetAllByRestaurant(ctx *gofr.Context, restaurantID int) ([]model.Rating, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, order_id, restaurant_id, rating, comment, created_at FROM order_ratings WHERE restaurant_id = ? ORDER BY created_at DESC",
		restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Rating
	for rows.Next() {
		var r model.Rating
		if err := rows.Scan(&r.ID, &r.OrderID, &r.RestaurantID, &r.Rating, &r.Comment, &r.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	if list == nil {
		list = []model.Rating{}
	}

	return list, nil
}

func (s *Rating) Create(ctx *gofr.Context, r *model.Rating) (*model.Rating, error) {
	now := time.Now()

	result, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO order_ratings (order_id, restaurant_id, rating, comment, created_at) VALUES (?, ?, ?, ?, ?)",
		r.OrderID, r.RestaurantID, r.Rating, r.Comment, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	r.ID = int(id)
	r.CreatedAt = now

	return r, nil
}
