package store

import (
	"qr-dinein-backend/model"
	"time"

	"gofr.dev/pkg/gofr"
)

type Product struct{}

func NewProduct() *Product {
	return &Product{}
}

func (s *Product) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Product, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, category_id, name, description, price, image, veg, available, prep_time, created_at, updated_at FROM products WHERE restaurant_id = ? ORDER BY id ASC",
		restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows)
}

func (s *Product) GetByCategory(ctx *gofr.Context, restaurantID, categoryID int) ([]model.Product, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, category_id, name, description, price, image, veg, available, prep_time, created_at, updated_at FROM products WHERE restaurant_id = ? AND category_id = ? ORDER BY id ASC",
		restaurantID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows)
}

func (s *Product) GetByID(ctx *gofr.Context, restaurantID, id int) (*model.Product, error) {
	var p model.Product
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, restaurant_id, category_id, name, description, price, image, veg, available, prep_time, created_at, updated_at FROM products WHERE id = ? AND restaurant_id = ?",
		id, restaurantID).
		Scan(&p.ID, &p.RestaurantID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Image, &p.Veg, &p.Available, &p.PrepTime, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *Product) Create(ctx *gofr.Context, p *model.Product) (*model.Product, error) {
	now := time.Now()

	result, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO products (restaurant_id, category_id, name, description, price, image, veg, available, prep_time, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		p.RestaurantID, p.CategoryID, p.Name, p.Description, p.Price, p.Image, p.Veg, p.Available, p.PrepTime, now, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	p.CreatedAt = now
	p.UpdatedAt = now

	return p, nil
}

func (s *Product) Update(ctx *gofr.Context, restaurantID, id int, p *model.Product) (*model.Product, error) {
	now := time.Now()

	_, err := ctx.SQL.ExecContext(ctx,
		"UPDATE products SET category_id = ?, name = ?, description = ?, price = ?, image = ?, veg = ?, available = ?, prep_time = ?, updated_at = ? WHERE id = ? AND restaurant_id = ?",
		p.CategoryID, p.Name, p.Description, p.Price, p.Image, p.Veg, p.Available, p.PrepTime, now, id, restaurantID)
	if err != nil {
		return nil, err
	}

	p.ID = id
	p.RestaurantID = restaurantID
	p.UpdatedAt = now

	return p, nil
}

func (s *Product) GetPrepTimes(ctx *gofr.Context, restaurantID int, productIDs []int) (map[int]int, error) {
	if len(productIDs) == 0 {
		return map[int]int{}, nil
	}

	placeholders := ""
	args := []interface{}{restaurantID}
	for i, id := range productIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}

	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, prep_time FROM products WHERE restaurant_id = ? AND id IN ("+placeholders+")",
		args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var id, prepTime int
		if err := rows.Scan(&id, &prepTime); err != nil {
			return nil, err
		}
		result[id] = prepTime
	}

	return result, nil
}

func (s *Product) Delete(ctx *gofr.Context, restaurantID, id int) error {
	_, err := ctx.SQL.ExecContext(ctx, "DELETE FROM products WHERE id = ? AND restaurant_id = ?", id, restaurantID)
	return err
}

type productRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

func scanProducts(rows productRows) ([]model.Product, error) {
	var list []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.RestaurantID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Image, &p.Veg, &p.Available, &p.PrepTime, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	if list == nil {
		list = []model.Product{}
	}

	return list, nil
}
