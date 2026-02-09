package store

import (
	"database/sql"
	"encoding/json"
	"qr-dinein-backend/model"
	"time"

	"gofr.dev/pkg/gofr"
)

type Order struct{}

func NewOrder() *Order {
	return &Order{}
}

func (s *Order) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Order, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, table_number, customer_mobile, customer_name, items, status, special_instructions, total, assigned_chef_id, estimated_ready_at, created_at, updated_at FROM orders WHERE restaurant_id = ? ORDER BY created_at DESC",
		restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (s *Order) GetByStatus(ctx *gofr.Context, restaurantID int, status string) ([]model.Order, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, table_number, customer_mobile, customer_name, items, status, special_instructions, total, assigned_chef_id, estimated_ready_at, created_at, updated_at FROM orders WHERE restaurant_id = ? AND status = ? ORDER BY created_at DESC",
		restaurantID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (s *Order) GetByPhone(ctx *gofr.Context, restaurantID int, phone string) ([]model.Order, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT id, restaurant_id, table_number, customer_mobile, customer_name, items, status, special_instructions, total, assigned_chef_id, estimated_ready_at, created_at, updated_at FROM orders WHERE restaurant_id = ? AND customer_mobile = ? ORDER BY created_at DESC",
		restaurantID, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (s *Order) GetByID(ctx *gofr.Context, restaurantID, id int) (*model.Order, error) {
	var o model.Order
	var itemsJSON []byte
	var tableNumber sql.NullString
	var chefID sql.NullInt64
	var estimatedReadyAt sql.NullTime

	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT id, restaurant_id, table_number, customer_mobile, customer_name, items, status, special_instructions, total, assigned_chef_id, estimated_ready_at, created_at, updated_at FROM orders WHERE id = ? AND restaurant_id = ?",
		id, restaurantID).
		Scan(&o.ID, &o.RestaurantID, &tableNumber, &o.CustomerMobile, &o.CustomerName, &itemsJSON, &o.Status, &o.SpecialInstructions, &o.Total, &chefID, &estimatedReadyAt, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if tableNumber.Valid {
		o.TableNumber = &tableNumber.String
	}

	if chefID.Valid {
		id := int(chefID.Int64)
		o.AssignedChefID = &id
	}

	if estimatedReadyAt.Valid {
		o.EstimatedReadyAt = &estimatedReadyAt.Time
	}

	if err := json.Unmarshal(itemsJSON, &o.Items); err != nil {
		return nil, err
	}

	return &o, nil
}

func (s *Order) Create(ctx *gofr.Context, o *model.Order) (*model.Order, error) {
	now := time.Now()

	itemsJSON, err := json.Marshal(o.Items)
	if err != nil {
		return nil, err
	}

	result, err := ctx.SQL.ExecContext(ctx,
		"INSERT INTO orders (restaurant_id, table_number, customer_mobile, customer_name, items, status, special_instructions, total, assigned_chef_id, estimated_ready_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		o.RestaurantID, o.TableNumber, o.CustomerMobile, o.CustomerName, string(itemsJSON), o.Status, o.SpecialInstructions, o.Total, o.AssignedChefID, o.EstimatedReadyAt, now, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	o.ID = int(id)
	o.CreatedAt = now
	o.UpdatedAt = now

	return o, nil
}

func (s *Order) Update(ctx *gofr.Context, restaurantID, id int, setClauses []string, args []interface{}) (*model.Order, error) {
	now := time.Now()
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	query := "UPDATE orders SET " + joinClauses(setClauses) + " WHERE id = ? AND restaurant_id = ?"
	args = append(args, id, restaurantID)

	_, err := ctx.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, restaurantID, id)
}

func joinClauses(clauses []string) string {
	result := ""
	for i, c := range clauses {
		if i > 0 {
			result += ", "
		}
		result += c
	}
	return result
}

func (s *Order) GetChefLoads(ctx *gofr.Context, restaurantID int) ([]model.ChefLoad, error) {
	rows, err := ctx.SQL.QueryContext(ctx,
		"SELECT assigned_chef_id, COUNT(*) FROM orders WHERE restaurant_id = ? AND status IN ('pending','preparing') AND assigned_chef_id IS NOT NULL GROUP BY assigned_chef_id",
		restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loads []model.ChefLoad
	for rows.Next() {
		var cl model.ChefLoad
		if err := rows.Scan(&cl.ChefID, &cl.OrderCount); err != nil {
			return nil, err
		}
		loads = append(loads, cl)
	}

	return loads, nil
}

func (s *Order) GetLeastRecentlyAssignedChef(ctx *gofr.Context, restaurantID int, chefIDs []int) (*int, error) {
	if len(chefIDs) == 0 {
		return nil, nil
	}

	placeholders := ""
	args := []interface{}{restaurantID}
	for i, id := range chefIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}

	query := `SELECT s.id FROM staff s
LEFT JOIN (
  SELECT assigned_chef_id, MAX(created_at) as last_assigned
  FROM orders WHERE restaurant_id = ? AND assigned_chef_id IS NOT NULL
  GROUP BY assigned_chef_id
) o ON s.id = o.assigned_chef_id
WHERE s.id IN (` + placeholders + `)
ORDER BY o.last_assigned ASC, s.id ASC LIMIT 1`

	var chefID int
	err := ctx.SQL.QueryRowContext(ctx, query, args...).Scan(&chefID)
	if err != nil {
		return nil, err
	}

	return &chefID, nil
}

func (s *Order) GetChefActiveOrderCount(ctx *gofr.Context, chefID int) (int, error) {
	var count int
	err := ctx.SQL.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM orders WHERE assigned_chef_id = ? AND status IN ('pending','preparing')",
		chefID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Order) Delete(ctx *gofr.Context, restaurantID, id int) error {
	_, err := ctx.SQL.ExecContext(ctx, "DELETE FROM orders WHERE id = ? AND restaurant_id = ?", id, restaurantID)
	return err
}

type orderRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

func scanOrders(rows orderRows) ([]model.Order, error) {
	var list []model.Order

	for rows.Next() {
		var o model.Order
		var itemsJSON []byte
		var tableNumber sql.NullString
		var chefID sql.NullInt64
		var estimatedReadyAt sql.NullTime

		if err := rows.Scan(&o.ID, &o.RestaurantID, &tableNumber, &o.CustomerMobile, &o.CustomerName, &itemsJSON, &o.Status, &o.SpecialInstructions, &o.Total, &chefID, &estimatedReadyAt, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}

		if tableNumber.Valid {
			o.TableNumber = &tableNumber.String
		}

		if chefID.Valid {
			id := int(chefID.Int64)
			o.AssignedChefID = &id
		}

		if estimatedReadyAt.Valid {
			o.EstimatedReadyAt = &estimatedReadyAt.Time
		}

		if err := json.Unmarshal(itemsJSON, &o.Items); err != nil {
			return nil, err
		}

		list = append(list, o)
	}

	if list == nil {
		list = []model.Order{}
	}

	return list, nil
}
