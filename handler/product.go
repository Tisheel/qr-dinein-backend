package handler

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/service"
	"strconv"

	"gofr.dev/pkg/gofr"
)

type Product struct {
	service *service.Product
}

func NewProduct(svc *service.Product) *Product {
	return &Product{service: svc}
}

func (h *Product) GetAll(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	// Filter by category if query param provided
	categoryIDStr := ctx.Param("categoryId")
	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid categoryId")
		}

		return h.service.GetByCategory(ctx, restaurantID, categoryID)
	}

	return h.service.GetAll(ctx, restaurantID)
}

func (h *Product) GetByID(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid product id")
	}

	return h.service.GetByID(ctx, restaurantID, id)
}

func (h *Product) Create(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	var p model.Product
	if err := ctx.Bind(&p); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Create(ctx, restaurantID, &p)
}

func (h *Product) Update(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid product id")
	}

	var p model.Product
	if err := ctx.Bind(&p); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Update(ctx, restaurantID, id, &p)
}

func (h *Product) Delete(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid product id")
	}

	if err := h.service.Delete(ctx, restaurantID, id); err != nil {
		return nil, err
	}

	return map[string]string{"message": "product deleted"}, nil
}
