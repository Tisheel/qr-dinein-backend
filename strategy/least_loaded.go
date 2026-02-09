package strategy

import (
	"fmt"
	"qr-dinein-backend/store"

	"gofr.dev/pkg/gofr"
)

type LeastLoadedStrategy struct {
	staffStore *store.Staff
	orderStore *store.Order
}

func NewLeastLoaded(staffStore *store.Staff, orderStore *store.Order) *LeastLoadedStrategy {
	return &LeastLoadedStrategy{staffStore: staffStore, orderStore: orderStore}
}

func (s *LeastLoadedStrategy) Assign(ctx *gofr.Context, restaurantID int) (*int, error) {
	chefs, err := s.staffStore.GetActiveChefs(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active chefs: %w", err)
	}

	if len(chefs) == 0 {
		return nil, nil
	}

	loads, err := s.orderStore.GetChefLoads(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chef loads: %w", err)
	}

	// Build load map
	loadMap := make(map[int]int)
	for _, l := range loads {
		loadMap[l.ChefID] = l.OrderCount
	}

	// Find chef with minimum load
	minLoad := -1
	var minChefID int

	for _, chef := range chefs {
		load := loadMap[chef.ID] // defaults to 0 if not in map
		if minLoad == -1 || load < minLoad {
			minLoad = load
			minChefID = chef.ID
		}
	}

	return &minChefID, nil
}
