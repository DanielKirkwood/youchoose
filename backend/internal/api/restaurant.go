package api

import (
	"context"
	"net/http"

	"github.com/DanielKirkwood/youchoose/internal/db"
	"github.com/DanielKirkwood/youchoose/internal/services"
	"github.com/danielgtaylor/huma/v2"
)

type GetNearestRestaurantsInput struct {
	Limit     int32   `query:"limit" required:"false" default:"30" doc:"the maximum number of rows to return" example:"100"`
	Offset    int32   `query:"offset" required:"false" default:"0" doc:"the row to start from" example:"20"`
	Radius    float64 `query:"radius" required:"false" default:"2" doc:"the radius, in miles, to search around the given location" example:"5" minimum:"1"`
	Longitude float64 `query:"longitude" required:"true" doc:"the longitude" example:"-4.251806" minimum:"-180.00000" maximum:"180.00000"`
	Latitude  float64 `query:"latitude" required:"true" doc:"the latitude" example:"55.864239" minimum:"-90.00000" maximum:"90.00000"`
}

type GetNearestRestaurantsBody struct {
	Restaurants []db.GetNearestRestaurantsRow `json:"restaurants"`
}

type GetNearestRestaurantsOutput struct {
	Body GetNearestRestaurantsBody
}

func RegisterRestaurantRoutes(api huma.API, restaurantService *services.RestaurantService) {
	huma.Register(api, huma.Operation{
		OperationID: "get-nearest-restaurants",
		Method:      http.MethodGet,
		Path:        "/restaurants/nearby",
		Summary:     "Get the restaurants within the given radius to the given location",
	}, func(ctx context.Context, i *GetNearestRestaurantsInput) (*GetNearestRestaurantsOutput, error) {
		radiusMetres := i.Radius * 1609.34

		params := db.GetNearestRestaurantsParams{
			Limit:         i.Limit,
			Offset:        i.Offset,
			MaxRadius:     radiusMetres,
			UserLongitude: i.Longitude,
			UserLatitude:  i.Latitude,
		}

		result, err := restaurantService.GetNearestRestaurants(ctx, params)
		if err != nil {
			return nil, huma.Error500InternalServerError("unable to fetch nearby restaurants", err)
		}

		return &GetNearestRestaurantsOutput{
			Body: GetNearestRestaurantsBody{
				Restaurants: result,
			},
		}, nil

	})
}
