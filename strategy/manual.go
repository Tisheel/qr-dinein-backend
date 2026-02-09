package strategy

import "gofr.dev/pkg/gofr"

type ManualStrategy struct{}

func (s *ManualStrategy) Assign(_ *gofr.Context, _ int) (*int, error) {
	return nil, nil
}
