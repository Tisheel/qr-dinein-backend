package store

import (
	"qr-dinein-backend/model"
	"time"

	"gofr.dev/pkg/gofr"
)

type Staff struct{}

func NewStaff() *Staff {
	return &Staff{}
}

func (s *Staff) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Staff, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, username, role, active, created_at, updated_at FROM staff WHERE restaurant_id = ? ORDER BY username ASC",
		restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Staff
	for rows.Next() {
		var st model.Staff
		if err := rows.Scan(&st.ID, &st.RestaurantID, &st.Username, &st.Role, &st.Active, &st.CreatedAt, &st.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, st)
	}

	if list == nil {
		list = []model.Staff{}
	}

	return list, nil
}

func (s *Staff) GetByID(ctx *gofr.Context, restaurantID, id int) (*model.Staff, error) {
	var st model.Staff
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, restaurant_id, username, role, active, created_at, updated_at FROM staff WHERE id = ? AND restaurant_id = ?",
		id, restaurantID).
		Scan(&st.ID, &st.RestaurantID, &st.Username, &st.Role, &st.Active, &st.CreatedAt, &st.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *Staff) GetByUsername(ctx *gofr.Context, username string) (*model.Staff, error) {
	var st model.Staff
	var pin *string
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, restaurant_id, username, pin, role, active, created_at, updated_at FROM staff WHERE username = ?",
		username).
		Scan(&st.ID, &st.RestaurantID, &st.Username, &pin, &st.Role, &st.Active, &st.CreatedAt, &st.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if pin != nil {
		st.Pin = *pin
	}

	return &st, nil
}

func (s *Staff) Create(ctx *gofr.Context, st *model.Staff) (*model.Staff, error) {
	now := time.Now()

	result, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO staff (restaurant_id, username, pin, role, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		st.RestaurantID, st.Username, st.Pin, st.Role, st.Active, now, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	st.ID = int(id)
	st.CreatedAt = now
	st.UpdatedAt = now
	st.Pin = "" // Don't return pin

	return st, nil
}

func (s *Staff) Update(ctx *gofr.Context, restaurantID, id int, st *model.Staff) (*model.Staff, error) {
	now := time.Now()

	if st.Pin != "" {
		_, err := ctx.SQL.ExecContext(ctx,
			"UPDATE staff SET username = ?, pin = ?, role = ?, active = ?, updated_at = ? WHERE id = ? AND restaurant_id = ?",
			st.Username, st.Pin, st.Role, st.Active, now, id, restaurantID)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := ctx.SQL.ExecContext(ctx,
			"UPDATE staff SET username = ?, role = ?, active = ?, updated_at = ? WHERE id = ? AND restaurant_id = ?",
			st.Username, st.Role, st.Active, now, id, restaurantID)
		if err != nil {
			return nil, err
		}
	}

	st.ID = id
	st.RestaurantID = restaurantID
	st.UpdatedAt = now
	st.Pin = "" // Don't return pin

	return st, nil
}

func (s *Staff) GetActiveChefs(ctx *gofr.Context, restaurantID int) ([]model.Staff, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, username, role, active, created_at, updated_at FROM staff WHERE restaurant_id = ? AND role = 'chef' AND active = TRUE ORDER BY id ASC",
		restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Staff
	for rows.Next() {
		var st model.Staff
		if err := rows.Scan(&st.ID, &st.RestaurantID, &st.Username, &st.Role, &st.Active, &st.CreatedAt, &st.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, st)
	}

	return list, nil
}

func (s *Staff) Delete(ctx *gofr.Context, restaurantID, id int) error {
	_, err := ctx.SQL.ExecContext(ctx, "DELETE FROM staff WHERE id = ? AND restaurant_id = ?", id, restaurantID)
	return err
}
