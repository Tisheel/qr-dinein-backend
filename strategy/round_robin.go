package strategy

import (
	"fmt"
	"qr-dinein-backend/store"

	"gofr.dev/pkg/gofr"
)

type RoundRobinStrategy struct {
	staffStore *store.Staff
}

func NewRoundRobin(staffStore *store.Staff) *RoundRobinStrategy {
	return &RoundRobinStrategy{staffStore: staffStore}
}

func (s *RoundRobinStrategy) Assign(ctx *gofr.Context, restaurantID int) (*int, error) {
	chefs, err := s.staffStore.GetActiveChefs(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active chefs: %w", err)
	}

	if len(chefs) == 0 {
		return nil, nil
	}

	key := fmt.Sprintf("chef_rr:%d", restaurantID)

	counter, err := ctx.Redis.Incr(ctx, key).Result()
	if err != nil {
		// Fallback to first chef if Redis fails
		id := chefs[0].ID
		return &id, nil
	}

	index := int((counter - 1) % int64(len(chefs)))
	id := chefs[index].ID

	return &id, nil
}
