package store

import (
	"qr-dinein-backend/model"

	"gofr.dev/pkg/gofr"
)

type Settings struct{}

func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Setting, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, `key`, value FROM settings WHERE restaurant_id = ? ORDER BY `key` ASC",
		restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Setting
	for rows.Next() {
		var st model.Setting
		if err := rows.Scan(&st.ID, &st.RestaurantID, &st.Key, &st.Value); err != nil {
			return nil, err
		}
		list = append(list, st)
	}

	if list == nil {
		list = []model.Setting{}
	}

	return list, nil
}

func (s *Settings) GetByKey(ctx *gofr.Context, restaurantID int, key string) (*model.Setting, error) {
	var st model.Setting
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, restaurant_id, `key`, value FROM settings WHERE restaurant_id = ? AND `key` = ?",
		restaurantID, key).
		Scan(&st.ID, &st.RestaurantID, &st.Key, &st.Value)
	if err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *Settings) Upsert(ctx *gofr.Context, restaurantID int, key, value string) (*model.Setting, error) {
	_, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO settings (restaurant_id, `key`, value) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE value = VALUES(value)",
		restaurantID, key, value)
	if err != nil {
		return nil, err
	}

	return s.GetByKey(ctx, restaurantID, key)
}

func (s *Settings) BulkUpsert(ctx *gofr.Context, restaurantID int, settings map[string]string) ([]model.Setting, error) {
	for key, value := range settings {
		_, err := ctx.SQL.ExecContext(ctx,
			"INSERT INTO settings (restaurant_id, `key`, value) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE value = VALUES(value)",
			restaurantID, key, value)
		if err != nil {
			return nil, err
		}
	}

	return s.GetAll(ctx, restaurantID)
}

func (s *Settings) Delete(ctx *gofr.Context, restaurantID int, key string) error {
	_, err := ctx.SQL.ExecContext(ctx, "DELETE FROM settings WHERE restaurant_id = ? AND `key` = ?", restaurantID, key)
	return err
}
