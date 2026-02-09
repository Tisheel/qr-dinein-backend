package service

import (
	"encoding/json"
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/store"
	"strconv"
	"strings"

	"gofr.dev/pkg/gofr"
)

type Category struct {
	store *store.Category
}

func NewCategory(s *store.Category) *Category {
	return &Category{store: s}
}

func (svc *Category) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Category, error) {
	cacheKey := "categories:" + strconv.Itoa(restaurantID)

	cached, err := ctx.Redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var categories []model.Category
		if err := json.Unmarshal([]byte(cached), &categories); err == nil {
			return categories, nil
		}
	}

	categories, err := svc.store.GetAll(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(categories); err == nil {
		ctx.Redis.Set(ctx, cacheKey, string(data), 0)
	}

	return categories, nil
}

func (svc *Category) GetByID(ctx *gofr.Context, restaurantID, id int) (*model.Category, error) {
	return svc.store.GetByID(ctx, restaurantID, id)
}

func (svc *Category) Create(ctx *gofr.Context, restaurantID int, c *model.Category) (*model.Category, error) {
	if strings.TrimSpace(c.Name) == "" {
		return nil, fmt.Errorf("category name is required")
	}

	c.RestaurantID = restaurantID

	result, err := svc.store.Create(ctx, c)
	if err != nil {
		return nil, err
	}

	svc.invalidateCache(ctx, restaurantID)

	return result, nil
}

func (svc *Category) Update(ctx *gofr.Context, restaurantID, id int, c *model.Category) (*model.Category, error) {
	result, err := svc.store.Update(ctx, restaurantID, id, c)
	if err != nil {
		return nil, err
	}

	svc.invalidateCache(ctx, restaurantID)

	return result, nil
}

func (svc *Category) Delete(ctx *gofr.Context, restaurantID, id int) error {
	if err := svc.store.Delete(ctx, restaurantID, id); err != nil {
		return err
	}

	svc.invalidateCache(ctx, restaurantID)

	return nil
}

func (svc *Category) invalidateCache(ctx *gofr.Context, restaurantID int) {
	cacheKey := "categories:" + strconv.Itoa(restaurantID)
	ctx.Redis.Del(ctx, cacheKey)
}
