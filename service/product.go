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

type Product struct {
	store *store.Product
}

func NewProduct(s *store.Product) *Product {
	return &Product{store: s}
}

func (svc *Product) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Product, error) {
	cacheKey := "products:" + strconv.Itoa(restaurantID)

	cached, err := ctx.Redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var products []model.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	products, err := svc.store.GetAll(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(products); err == nil {
		ctx.Redis.Set(ctx, cacheKey, string(data), 0)
	}

	return products, nil
}

func (svc *Product) GetByCategory(ctx *gofr.Context, restaurantID, categoryID int) ([]model.Product, error) {
	return svc.store.GetByCategory(ctx, restaurantID, categoryID)
}

func (svc *Product) GetByID(ctx *gofr.Context, restaurantID, id int) (*model.Product, error) {
	return svc.store.GetByID(ctx, restaurantID, id)
}

func (svc *Product) Create(ctx *gofr.Context, restaurantID int, p *model.Product) (*model.Product, error) {
	if strings.TrimSpace(p.Name) == "" {
		return nil, fmt.Errorf("product name is required")
	}

	if p.Price <= 0 {
		return nil, fmt.Errorf("product price must be greater than 0")
	}

	if p.CategoryID == 0 {
		return nil, fmt.Errorf("category id is required")
	}

	p.RestaurantID = restaurantID

	result, err := svc.store.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	svc.invalidateCache(ctx, restaurantID)

	return result, nil
}

func (svc *Product) Update(ctx *gofr.Context, restaurantID, id int, p *model.Product) (*model.Product, error) {
	result, err := svc.store.Update(ctx, restaurantID, id, p)
	if err != nil {
		return nil, err
	}

	svc.invalidateCache(ctx, restaurantID)

	return result, nil
}

func (svc *Product) Delete(ctx *gofr.Context, restaurantID, id int) error {
	if err := svc.store.Delete(ctx, restaurantID, id); err != nil {
		return err
	}

	svc.invalidateCache(ctx, restaurantID)

	return nil
}

func (svc *Product) invalidateCache(ctx *gofr.Context, restaurantID int) {
	cacheKey := "products:" + strconv.Itoa(restaurantID)
	ctx.Redis.Del(ctx, cacheKey)
}
