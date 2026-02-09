package handler

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/service"
	"strconv"

	"gofr.dev/pkg/gofr"
)

type Staff struct {
	service *service.Staff
}

func NewStaff(svc *service.Staff) *Staff {
	return &Staff{service: svc}
}

func (h *Staff) GetAll(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	return h.service.GetAll(ctx, restaurantID)
}

func (h *Staff) GetByID(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid staff id")
	}

	return h.service.GetByID(ctx, restaurantID, id)
}

func (h *Staff) Create(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	var s model.Staff
	if err := ctx.Bind(&s); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Create(ctx, restaurantID, &s)
}

func (h *Staff) Update(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid staff id")
	}

	var s model.Staff
	if err := ctx.Bind(&s); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Update(ctx, restaurantID, id, &s)
}

func (h *Staff) Delete(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid staff id")
	}

	if err := h.service.Delete(ctx, restaurantID, id); err != nil {
		return nil, err
	}

	return map[string]string{"message": "staff member deleted"}, nil
}
