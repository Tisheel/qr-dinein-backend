package strategy

import (
	"qr-dinein-backend/store"

	"gofr.dev/pkg/gofr"
)

type Resolver struct {
	settingsStore *store.Settings
	staffStore    *store.Staff
	orderStore    *store.Order
}

func NewResolver(settingsStore *store.Settings, staffStore *store.Staff, orderStore *store.Order) *Resolver {
	return &Resolver{
		settingsStore: settingsStore,
		staffStore:    staffStore,
		orderStore:    orderStore,
	}
}

func (r *Resolver) Resolve(ctx *gofr.Context, restaurantID int) ChefAssigner {
	setting, err := r.settingsStore.GetByKey(ctx, restaurantID, "chef_assignment_strategy")
	if err != nil {
		return &ManualStrategy{}
	}

	switch setting.Value {
	case StrategyRoundRobin:
		return NewRoundRobin(r.staffStore)
	case StrategyLeastLoaded:
		return NewLeastLoaded(r.staffStore, r.orderStore)
	case StrategyRandom:
		return NewRandom(r.staffStore)
	case StrategyLeastRecentlyAssigned:
		return NewLeastRecentlyAssigned(r.staffStore, r.orderStore)
	default:
		return &ManualStrategy{}
	}
}
