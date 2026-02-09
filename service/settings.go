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

type Settings struct {
	store *store.Settings
}

func NewSettings(s *store.Settings) *Settings {
	return &Settings{store: s}
}

func (svc *Settings) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Setting, error) {
	cacheKey := "settings:" + strconv.Itoa(restaurantID)

	cached, err := ctx.Redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var settings []model.Setting
		if err := json.Unmarshal([]byte(cached), &settings); err == nil {
			return settings, nil
		}
	}

	settings, err := svc.store.GetAll(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(settings); err == nil {
		ctx.Redis.Set(ctx, cacheKey, string(data), 0)
	}

	return settings, nil
}

func (svc *Settings) GetByKey(ctx *gofr.Context, restaurantID int, key string) (*model.Setting, error) {
	return svc.store.GetByKey(ctx, restaurantID, key)
}

func (svc *Settings) Upsert(ctx *gofr.Context, restaurantID int, key, value string) (*model.Setting, error) {
	if strings.TrimSpace(key) == "" {
		return nil, fmt.Errorf("setting key is required")
	}

	result, err := svc.store.Upsert(ctx, restaurantID, key, value)
	if err != nil {
		return nil, err
	}

	svc.invalidateCache(ctx, restaurantID)

	return result, nil
}

func (svc *Settings) BulkUpsert(ctx *gofr.Context, restaurantID int, settings map[string]string) ([]model.Setting, error) {
	result, err := svc.store.BulkUpsert(ctx, restaurantID, settings)
	if err != nil {
		return nil, err
	}

	svc.invalidateCache(ctx, restaurantID)

	return result, nil
}

func (svc *Settings) Delete(ctx *gofr.Context, restaurantID int, key string) error {
	if err := svc.store.Delete(ctx, restaurantID, key); err != nil {
		return err
	}

	svc.invalidateCache(ctx, restaurantID)

	return nil
}

func (svc *Settings) invalidateCache(ctx *gofr.Context, restaurantID int) {
	cacheKey := "settings:" + strconv.Itoa(restaurantID)
	ctx.Redis.Del(ctx, cacheKey)
}
