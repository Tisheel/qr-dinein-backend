package handler

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/service"
	"strconv"

	"gofr.dev/pkg/gofr"
)

type Category struct {
	service *service.Category
}

func NewCategory(svc *service.Category) *Category {
	return &Category{service: svc}
}

func (h *Category) GetAll(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	return h.service.GetAll(ctx, restaurantID)
}

func (h *Category) GetByID(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid category id")
	}

	return h.service.GetByID(ctx, restaurantID, id)
}

func (h *Category) Create(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	var c model.Category
	if err := ctx.Bind(&c); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Create(ctx, restaurantID, &c)
}

func (h *Category) Update(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid category id")
	}

	var c model.Category
	if err := ctx.Bind(&c); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Update(ctx, restaurantID, id, &c)
}

func (h *Category) Delete(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid category id")
	}

	if err := h.service.Delete(ctx, restaurantID, id); err != nil {
		return nil, err
	}

	return map[string]string{"message": "category deleted"}, nil
}
