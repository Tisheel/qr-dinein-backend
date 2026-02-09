package strategy

import (
	"fmt"
	"math/rand"
	"qr-dinein-backend/store"

	"gofr.dev/pkg/gofr"
)

type RandomStrategy struct {
	staffStore *store.Staff
}

func NewRandom(staffStore *store.Staff) *RandomStrategy {
	return &RandomStrategy{staffStore: staffStore}
}

func (s *RandomStrategy) Assign(ctx *gofr.Context, restaurantID int) (*int, error) {
	chefs, err := s.staffStore.GetActiveChefs(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active chefs: %w", err)
	}

	if len(chefs) == 0 {
		return nil, nil
	}

	index := rand.Intn(len(chefs))
	id := chefs[index].ID

	return &id, nil
}
