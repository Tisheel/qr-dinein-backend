package store

import (
	"qr-dinein-backend/model"
	"time"

	"gofr.dev/pkg/gofr"
)

type Restaurant struct{}

func NewRestaurant() *Restaurant {
	return &Restaurant{}
}

func (s *Restaurant) GetAll(ctx *gofr.Context) ([]model.Restaurant, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, name, slug, address, phone, logo, currency, tax_rate, active, created_at, updated_at FROM restaurants ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		if err := rows.Scan(&r.ID, &r.Name, &r.Slug, &r.Address, &r.Phone, &r.Logo, &r.Currency, &r.TaxRate, &r.Active, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	if list == nil {
		list = []model.Restaurant{}
	}

	return list, nil
}

func (s *Restaurant) GetByID(ctx *gofr.Context, id int) (*model.Restaurant, error) {
	var r model.Restaurant
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, name, slug, address, phone, logo, currency, tax_rate, active, created_at, updated_at FROM restaurants WHERE id = ?", id).
		Scan(&r.ID, &r.Name, &r.Slug, &r.Address, &r.Phone, &r.Logo, &r.Currency, &r.TaxRate, &r.Active, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *Restaurant) GetBySlug(ctx *gofr.Context, slug string) (*model.Restaurant, error) {
	var r model.Restaurant
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, name, slug, address, phone, logo, currency, tax_rate, active, created_at, updated_at FROM restaurants WHERE slug = ?", slug).
		Scan(&r.ID, &r.Name, &r.Slug, &r.Address, &r.Phone, &r.Logo, &r.Currency, &r.TaxRate, &r.Active, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *Restaurant) Create(ctx *gofr.Context, r *model.Restaurant) (*model.Restaurant, error) {
	now := time.Now()

	result, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO restaurants (name, slug, address, phone, logo, currency, tax_rate, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		r.Name, r.Slug, r.Address, r.Phone, r.Logo, r.Currency, r.TaxRate, r.Active, now, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	r.ID = int(id)
	r.CreatedAt = now
	r.UpdatedAt = now

	return r, nil
}

func (s *Restaurant) Update(ctx *gofr.Context, id int, r *model.Restaurant) (*model.Restaurant, error) {
	now := time.Now()

	_, err := ctx.SQL.ExecContext(ctx,
		"UPDATE restaurants SET name = ?, slug = ?, address = ?, phone = ?, logo = ?, currency = ?, tax_rate = ?, active = ?, updated_at = ? WHERE id = ?",
		r.Name, r.Slug, r.Address, r.Phone, r.Logo, r.Currency, r.TaxRate, r.Active, now, id)
	if err != nil {
		return nil, err
	}

	r.ID = id
	r.UpdatedAt = now

	return r, nil
}

func (s *Restaurant) Delete(ctx *gofr.Context, id int) error {
	_, err := ctx.SQL.ExecContext(ctx, "DELETE FROM restaurants WHERE id = ?", id)
	return err
}
