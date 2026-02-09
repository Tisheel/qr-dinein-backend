package store

import (
	"qr-dinein-backend/model"
	"time"

	"gofr.dev/pkg/gofr"
)

type Category struct{}

func NewCategory() *Category {
	return &Category{}
}

func (s *Category) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Category, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, name, `order`, image, created_at, updated_at FROM categories WHERE restaurant_id = ? ORDER BY `order` ASC",
		restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.RestaurantID, &c.Name, &c.Order, &c.Image, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	if list == nil {
		list = []model.Category{}
	}

	return list, nil
}

func (s *Category) GetByID(ctx *gofr.Context, restaurantID, id int) (*model.Category, error) {
	var c model.Category
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, restaurant_id, name, `order`, image, created_at, updated_at FROM categories WHERE id = ? AND restaurant_id = ?",
		id, restaurantID).
		Scan(&c.ID, &c.RestaurantID, &c.Name, &c.Order, &c.Image, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *Category) Create(ctx *gofr.Context, c *model.Category) (*model.Category, error) {
	now := time.Now()

	result, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO categories (restaurant_id, name, `order`, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		c.RestaurantID, c.Name, c.Order, c.Image, now, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	c.ID = int(id)
	c.CreatedAt = now
	c.UpdatedAt = now

	return c, nil
}

func (s *Category) Update(ctx *gofr.Context, restaurantID, id int, c *model.Category) (*model.Category, error) {
	now := time.Now()

	_, err := ctx.SQL.ExecContext(ctx,
		"UPDATE categories SET name = ?, `order` = ?, image = ?, updated_at = ? WHERE id = ? AND restaurant_id = ?",
		c.Name, c.Order, c.Image, now, id, restaurantID)
	if err != nil {
		return nil, err
	}

	c.ID = id
	c.RestaurantID = restaurantID
	c.UpdatedAt = now

	return c, nil
}

func (s *Category) Delete(ctx *gofr.Context, restaurantID, id int) error {
	_, err := ctx.SQL.ExecContext(ctx, "DELETE FROM categories WHERE id = ? AND restaurant_id = ?", id, restaurantID)
	return err
}
