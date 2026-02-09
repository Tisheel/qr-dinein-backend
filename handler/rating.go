package handler

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/service"
	"strconv"

	"gofr.dev/pkg/gofr"
)

type Rating struct {
	service *service.RatingService
}

func NewRating(svc *service.RatingService) *Rating {
	return &Rating{service: svc}
}

func (h *Rating) GetByOrderID(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	orderID, err := strconv.Atoi(ctx.PathParam("orderId"))
	if err != nil {
		return nil, fmt.Errorf("invalid order id")
	}

	return h.service.GetByOrderID(ctx, restaurantID, orderID)
}

func (h *Rating) GetAllByRestaurant(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	return h.service.GetAllByRestaurant(ctx, restaurantID)
}

func (h *Rating) Create(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	orderID, err := strconv.Atoi(ctx.PathParam("orderId"))
	if err != nil {
		return nil, fmt.Errorf("invalid order id")
	}

	var r model.Rating
	if err := ctx.Bind(&r); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Create(ctx, restaurantID, orderID, &r)
}
