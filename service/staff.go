package service

import (
	"fmt"
	"qr-dinein-backend/model"
	"qr-dinein-backend/store"
	"strings"

	"gofr.dev/pkg/gofr"
)

type Staff struct {
	store *store.Staff
}

func NewStaff(s *store.Staff) *Staff {
	return &Staff{store: s}
}

func (svc *Staff) GetAll(ctx *gofr.Context, restaurantID int) ([]model.Staff, error) {
	return svc.store.GetAll(ctx, restaurantID)
}

func (svc *Staff) GetByID(ctx *gofr.Context, restaurantID, id int) (*model.Staff, error) {
	return svc.store.GetByID(ctx, restaurantID, id)
}

func (svc *Staff) Create(ctx *gofr.Context, restaurantID int, st *model.Staff) (*model.Staff, error) {
	if strings.TrimSpace(st.Username) == "" {
		return nil, fmt.Errorf("staff username is required")
	}

	if st.Role == "superuser" {
		return nil, fmt.Errorf("cannot assign superuser role to staff")
	}

	if st.Pin == "" {
		return nil, fmt.Errorf("valid pin is required")
	}

	if len(st.Pin) != 6 {
		return nil, fmt.Errorf("invalid pin length")
	}

	if st.Role == "" {
		st.Role = "chef"
	}

	st.RestaurantID = restaurantID
	st.Active = true

	return svc.store.Create(ctx, st)
}

func (svc *Staff) Update(ctx *gofr.Context, restaurantID, id int, st *model.Staff) (*model.Staff, error) {
	if st.Role == "superuser" {
		return nil, fmt.Errorf("cannot assign superuser role to staff")
	}

	if _, err := svc.store.GetByID(ctx, restaurantID, id); err != nil {
		return nil, fmt.Errorf("staff member not found: %w", err)
	}

	return svc.store.Update(ctx, restaurantID, id, st)
}

func (svc *Staff) Delete(ctx *gofr.Context, restaurantID, id int) error {
	if _, err := svc.store.GetByID(ctx, restaurantID, id); err != nil {
		return fmt.Errorf("staff member not found: %w", err)
	}

	return svc.store.Delete(ctx, restaurantID, id)
}
