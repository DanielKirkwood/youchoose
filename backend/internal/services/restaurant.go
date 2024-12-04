package services

import (
	"context"

	"github.com/DanielKirkwood/youchoose/internal/db"
)

type RestaurantService struct {
	Queries *db.Queries
}

func (s *RestaurantService) GetNearestRestaurants(ctx context.Context, arg db.GetNearestRestaurantsParams) ([]db.GetNearestRestaurantsRow, error) {
	return s.Queries.GetNearestRestaurants(ctx, arg)
}
