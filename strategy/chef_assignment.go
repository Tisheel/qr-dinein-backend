package strategy

import "gofr.dev/pkg/gofr"

const (
	StrategyManual               = "manual"
	StrategyRoundRobin           = "round_robin"
	StrategyLeastLoaded          = "least_loaded"
	StrategyRandom               = "random"
	StrategyLeastRecentlyAssigned = "least_recently_assigned"
)

// ChefAssigner assigns a chef to an order for a given restaurant.
// Returns a pointer to the chef ID, or nil if no chef is available.
type ChefAssigner interface {
	Assign(ctx *gofr.Context, restaurantID int) (*int, error)
}
