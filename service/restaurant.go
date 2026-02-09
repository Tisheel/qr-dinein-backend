package service

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/store"
	"strings"

	"gofr.dev/pkg/gofr"
)

type Restaurant struct {
	store *store.Restaurant
}

func NewRestaurant(s *store.Restaurant) *Restaurant {
	return &Restaurant{store: s}
}

func (svc *Restaurant) GetAll(ctx *gofr.Context) ([]model.Restaurant, error) {
	return svc.store.GetAll(ctx)
}

func (svc *Restaurant) GetByID(ctx *gofr.Context, id int) (*model.Restaurant, error) {
	return svc.store.GetByID(ctx, id)
}

func (svc *Restaurant) GetBySlug(ctx *gofr.Context, slug string) (*model.Restaurant, error) {
	return svc.store.GetBySlug(ctx, slug)
}

func (svc *Restaurant) Create(ctx *gofr.Context, r *model.Restaurant) (*model.Restaurant, error) {
	if strings.TrimSpace(r.Name) == "" {
		return nil, fmt.Errorf("restaurant name is required")
	}

	if strings.TrimSpace(r.Slug) == "" {
		r.Slug = generateSlug(r.Name)
	}

	if r.Currency == "" {
		r.Currency = "INR"
	}

	r.Active = true

	return svc.store.Create(ctx, r)
}

func (svc *Restaurant) Update(ctx *gofr.Context, id int, r *model.Restaurant) (*model.Restaurant, error) {
	if _, err := svc.store.GetByID(ctx, id); err != nil {
		return nil, fmt.Errorf("restaurant not found: %w", err)
	}

	return svc.store.Update(ctx, id, r)
}

func (svc *Restaurant) Delete(ctx *gofr.Context, id int) error {
	if _, err := svc.store.GetByID(ctx, id); err != nil {
		return fmt.Errorf("restaurant not found: %w", err)
	}

	return svc.store.Delete(ctx, id)
}

func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}
