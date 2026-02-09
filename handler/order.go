package handler

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/service"
	"strconv"

	"gofr.dev/pkg/gofr"
)

type Order struct {
	service *service.Order
}

func NewOrder(svc *service.Order) *Order {
	return &Order{service: svc}
}

func (h *Order) GetAll(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	// Filter by status if query param provided
	status := ctx.Param("status")
	if status != "" {
		return h.service.GetByStatus(ctx, restaurantID, status)
	}

	return h.service.GetAll(ctx, restaurantID)
}

func (h *Order) GetByID(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid order id")
	}

	return h.service.GetByID(ctx, restaurantID, id)
}

func (h *Order) GetByPhone(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	phone := ctx.Param("phone")
	if phone == "" {
		return nil, fmt.Errorf("phone query parameter is required")
	}

	return h.service.GetByPhone(ctx, restaurantID, phone)
}

func (h *Order) Create(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	var o model.Order
	if err := ctx.Bind(&o); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Create(ctx, restaurantID, &o)
}

func (h *Order) Update(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid order id")
	}

	var o model.Order
	if err := ctx.Bind(&o); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Update(ctx, restaurantID, id, &o)
}

func (h *Order) Delete(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid order id")
	}

	if err := h.service.Delete(ctx, restaurantID, id); err != nil {
		return nil, err
	}

	return map[string]string{"message": "order deleted"}, nil
}
