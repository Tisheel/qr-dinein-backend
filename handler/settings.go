package handler

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/service"
	"strconv"

	"gofr.dev/pkg/gofr"
)

type Settings struct {
	service *service.Settings
}

func NewSettings(svc *service.Settings) *Settings {
	return &Settings{service: svc}
}

func (h *Settings) GetAll(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	return h.service.GetAll(ctx, restaurantID)
}

func (h *Settings) GetByKey(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	key := ctx.PathParam("key")
	if key == "" {
		return nil, fmt.Errorf("setting key is required")
	}

	return h.service.GetByKey(ctx, restaurantID, key)
}

func (h *Settings) BulkUpsert(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	var settings map[string]string
	if err := ctx.Bind(&settings); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.BulkUpsert(ctx, restaurantID, settings)
}

func (h *Settings) Upsert(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	key := ctx.PathParam("key")
	if key == "" {
		return nil, fmt.Errorf("setting key is required")
	}

	var body model.Setting
	if err := ctx.Bind(&body); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return h.service.Upsert(ctx, restaurantID, key, body.Value)
}

func (h *Settings) Delete(ctx *gofr.Context) (interface{}, error) {
	restaurantID, err := strconv.Atoi(ctx.PathParam("restaurantId"))
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant id")
	}

	key := ctx.PathParam("key")
	if key == "" {
		return nil, fmt.Errorf("setting key is required")
	}

	if err := h.service.Delete(ctx, restaurantID, key); err != nil {
		return nil, err
	}

	return map[string]string{"message": "setting deleted"}, nil
}
