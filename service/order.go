package service

import (
	"encoding/json"
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/store"
	"qr-dinein-backend/strategy"
	"time"

	"gofr.dev/pkg/gofr"
)

type Order struct {
	store         *store.Order
	productStore  *store.Product
	settingsStore *store.Settings
	customerSvc   *Customer
	chefResolver  *strategy.Resolver
}

func NewOrder(s *store.Order, productStore *store.Product, settingsStore *store.Settings, customerSvc *Customer, chefResolver *strategy.Resolver) *Order {
	return &Order{
		store:         s,
		productStore:  productStore,
		settingsStore: settingsStore,
		customerSvc:   customerSvc,
		chefResolver:  chefResolver,
	}
}

func (svc *Order) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Order, error) {
	return svc.store.GetAll(ctx, restaurantID)
}

func (svc *Order) GetByStatus(ctx *gofr.Context, restaurantID int, status string) ([]model.Order, error) {
	return svc.store.GetByStatus(ctx, restaurantID, status)
}

func (svc *Order) GetByPhone(ctx *gofr.Context, restaurantID int, phone string) ([]model.Order, error) {
	return svc.store.GetByPhone(ctx, restaurantID, phone)
}

func (svc *Order) GetByID(ctx *gofr.Context, restaurantID, id int) (*model.Order, error) {
	return svc.store.GetByID(ctx, restaurantID, id)
}

func (svc *Order) Create(ctx *gofr.Context, restaurantID int, o *model.Order) (*model.Order, error) {
	if len(o.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	// Check if customer auth is required
	authRequired := false
	if setting, err := svc.settingsStore.GetByKey(ctx, restaurantID, "customer_auth_required"); err == nil && setting.Value == "true" {
		authRequired = true
	}

	if authRequired {
		// Customer must verify phone via OTP before placing an order
		if o.SessionToken == "" {
			return nil, fmt.Errorf("customer authentication is required: please verify your phone number first")
		}

		sessionToken := o.SessionToken

		session, err := svc.customerSvc.GetSession(ctx, sessionToken)
		if err != nil {
			return nil, fmt.Errorf("invalid or expired session: %w", err)
		}

		if session.RestaurantID != restaurantID {
			return nil, fmt.Errorf("session is not valid for this restaurant")
		}

		o.CustomerMobile = session.PhoneNumber

		defer func() {
			_ = svc.customerSvc.InvalidateSession(ctx, sessionToken)
		}()
	} else if o.SessionToken != "" {
		// Auth not required but session token provided — still validate it
		sessionToken := o.SessionToken

		session, err := svc.customerSvc.GetSession(ctx, sessionToken)
		if err != nil {
			return nil, fmt.Errorf("invalid or expired session: %w", err)
		}

		if session.RestaurantID != restaurantID {
			return nil, fmt.Errorf("session is not valid for this restaurant")
		}

		o.CustomerMobile = session.PhoneNumber

		defer func() {
			_ = svc.customerSvc.InvalidateSession(ctx, sessionToken)
		}()
	} else if o.CustomerMobile == "" {
		return nil, fmt.Errorf("customer mobile is required")
	}

	// Clear session token before storing (not persisted)
	o.SessionToken = ""

	// Calculate total from items
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}

	o.Total = total
	o.RestaurantID = restaurantID

	if o.Status == "" {
		o.Status = "pending"
	}

	// Auto-assign chef
	assigner := svc.chefResolver.Resolve(ctx, restaurantID)
	chefID, err := assigner.Assign(ctx, restaurantID)
	if err != nil {
		ctx.Logger.Errorf("chef auto-assignment failed: %v", err)
	} else {
		o.AssignedChefID = chefID
	}

	// Calculate estimated ready time
	svc.calculateEstimatedReadyAt(ctx, o)

	return svc.store.Create(ctx, o)
}

func (svc *Order) Update(ctx *gofr.Context, restaurantID, id int, o *model.Order) (*model.Order, error) {
	existing, err := svc.store.GetByID(ctx, restaurantID, id)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Chef can only be assigned when order is pending
	if o.AssignedChefID != nil && existing.Status != "pending" {
		return nil, fmt.Errorf("chef can only be assigned when order is in pending state")
	}

	// Validate status transitions
	if o.Status != "" && o.Status != existing.Status {
		if !isValidStatusTransition(existing.Status, o.Status) {
			return nil, fmt.Errorf("invalid status transition from '%s' to '%s'", existing.Status, o.Status)
		}
	}

	// Build partial update
	var setClauses []string
	var args []interface{}

	if o.Status != "" {
		setClauses = append(setClauses, "status = ?")
		args = append(args, o.Status)
	}
	if o.AssignedChefID != nil {
		setClauses = append(setClauses, "assigned_chef_id = ?")
		args = append(args, *o.AssignedChefID)
	}
	if o.TableNumber != nil {
		setClauses = append(setClauses, "table_number = ?")
		args = append(args, *o.TableNumber)
	}
	if o.CustomerMobile != "" {
		setClauses = append(setClauses, "customer_mobile = ?")
		args = append(args, o.CustomerMobile)
	}
	if o.CustomerName != "" {
		setClauses = append(setClauses, "customer_name = ?")
		args = append(args, o.CustomerName)
	}
	if len(o.Items) > 0 {
		setClauses = append(setClauses, "items = ?")
		itemsJSON, err := json.Marshal(o.Items)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal items: %w", err)
		}
		args = append(args, string(itemsJSON))

		// Recalculate total
		var total float64
		for _, item := range o.Items {
			total += item.Price * float64(item.Quantity)
		}
		setClauses = append(setClauses, "total = ?")
		args = append(args, total)
	}
	if o.SpecialInstructions != "" {
		setClauses = append(setClauses, "special_instructions = ?")
		args = append(args, o.SpecialInstructions)
	}

	if len(setClauses) == 0 {
		return existing, nil
	}

	return svc.store.Update(ctx, restaurantID, id, setClauses, args)
}

func (svc *Order) Delete(ctx *gofr.Context, restaurantID, id int) error {
	return svc.store.Delete(ctx, restaurantID, id)
}

const defaultPrepTime = 5 // minutes

func (svc *Order) calculateEstimatedReadyAt(ctx *gofr.Context, o *model.Order) {
	// Collect unique product IDs
	productIDs := make([]int, 0, len(o.Items))
	seen := make(map[int]bool)
	for _, item := range o.Items {
		if !seen[item.ProductID] {
			productIDs = append(productIDs, item.ProductID)
			seen[item.ProductID] = true
		}
	}

	// Look up prep times from DB
	prepTimes, err := svc.productStore.GetPrepTimes(ctx, o.RestaurantID, productIDs)
	if err != nil {
		ctx.Logger.Errorf("failed to get prep times: %v", err)
	}

	// Max prep time across all items (parallel prep)
	maxPrepTime := 0
	for _, item := range o.Items {
		pt := defaultPrepTime
		if t, ok := prepTimes[item.ProductID]; ok && t > 0 {
			pt = t
		}
		if pt > maxPrepTime {
			maxPrepTime = pt
		}
	}

	if maxPrepTime == 0 {
		maxPrepTime = defaultPrepTime
	}

	// Factor in chef's queue depth
	queueDepth := 0
	if o.AssignedChefID != nil {
		count, err := svc.store.GetChefActiveOrderCount(ctx, *o.AssignedChefID)
		if err != nil {
			ctx.Logger.Errorf("failed to get chef active order count: %v", err)
		} else {
			queueDepth = count
		}
	}

	estimatedMinutes := maxPrepTime + (queueDepth * maxPrepTime)
	readyAt := time.Now().Add(time.Duration(estimatedMinutes) * time.Minute)
	o.EstimatedReadyAt = &readyAt
}

func isValidStatusTransition(from, to string) bool {
	transitions := map[string][]string{
		"pending":   {"preparing", "cancelled"},
		"preparing": {"completed", "cancelled"},
		"completed": {},
		"cancelled": {},
	}

	allowed, ok := transitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}

	return false
}
