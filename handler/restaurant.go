package handler

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/service"
	"strconv"

	"gofr.dev/pkg/gofr"
)

type Restaurant struct {
	service *service.Restaurant
}

func NewRestaurant(svc *service.Restaurant) *Restaurant {
	return &Restaurant{service: svc}
}

func (h *Restaurant) GetAll(ctx *gofr.Context) (interface{}, error) {
	return h.service.GetAll(ctx)
}

func (h *Restaurant) GetByID(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	return h.service.GetByID(ctx, id)
}

func (h *Restaurant) GetBySlug(ctx *gofr.Context) (interface{}, error) {
	slug := ctx.PathParam("slug")
	if slug == "" {
		return nil, fmt.Errorf("slug is required")
	}

	return h.service.GetBySlug(ctx, slug)
}

func (h *Restaurant) Create(ctx *gofr.Context) (interface{}, error) {
	var r model.Restaurant
	if err := ctx.Bind(&r); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Create(ctx, &r)
}

func (h *Restaurant) Update(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	var r model.Restaurant
	if err := ctx.Bind(&r); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Update(ctx, id, &r)
}

func (h *Restaurant) Delete(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	if err := h.service.Delete(ctx, id); err != nil {
		return nil, err
	}

	return map[string]string{"message": "restaurant deleted"}, nil
}
