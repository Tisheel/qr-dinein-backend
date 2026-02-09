package strategy

import (
	"fmt"
	"qr-dinein-backend/store"

	"gofr.dev/pkg/gofr"
)

type LeastRecentlyAssignedStrategy struct {
	staffStore *store.Staff
	orderStore *store.Order
}

func NewLeastRecentlyAssigned(staffStore *store.Staff, orderStore *store.Order) *LeastRecentlyAssignedStrategy {
	return &LeastRecentlyAssignedStrategy{staffStore: staffStore, orderStore: orderStore}
}

func (s *LeastRecentlyAssignedStrategy) Assign(ctx *gofr.Context, restaurantID int) (*int, error) {
	chefs, err := s.staffStore.GetActiveChefs(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active chefs: %w", err)
	}

	if len(chefs) == 0 {
		return nil, nil
	}

	chefIDs := make([]int, len(chefs))
	for i, c := range chefs {
		chefIDs[i] = c.ID
	}

	return s.orderStore.GetLeastRecentlyAssignedChef(ctx, restaurantID, chefIDs)
}
