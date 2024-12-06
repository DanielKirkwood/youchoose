// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: restaurants.sql

package db

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const getNearestRestaurants = `-- name: GetNearestRestaurants :many
with user_location as (
    select st_setsrid(
            st_makepoint($5::float, $6::float),
            4326
        ) as location
)
select r.id,
    r.name,
    r.address_line1,
    r.address_line2,
    r.postcode,
    r.longitude,
    r.latitude,
    st_distance(ul.location, r.geolocation) as distance_meters
from restaurants r,
    user_location ul
where r.geolocation is not null
    and st_dwithin(ul.location, r.geolocation, $3::float)
    and (
        (
            r.valid = true
            and r.created_by is null
        )
        or r.created_by = $4::uuid
    )
order by distance_meters asc
limit $1 offset $2
`

type GetNearestRestaurantsParams struct {
	Limit         int32     `json:"limit"`
	Offset        int32     `json:"offset"`
	MaxRadius     float64   `json:"max_radius"`
	UserID        uuid.UUID `json:"user_id"`
	UserLongitude float64   `json:"user_longitude"`
	UserLatitude  float64   `json:"user_latitude"`
}

type GetNearestRestaurantsRow struct {
	ID             uuid.UUID     `json:"id"`
	Name           string        `json:"name"`
	AddressLine1   pgtype.Text   `json:"address_line1"`
	AddressLine2   pgtype.Text   `json:"address_line2"`
	Postcode       pgtype.Text   `json:"postcode"`
	Longitude      pgtype.Float8 `json:"longitude"`
	Latitude       pgtype.Float8 `json:"latitude"`
	DistanceMeters interface{}   `json:"distance_meters"`
}

func (q *Queries) GetNearestRestaurants(ctx context.Context, arg GetNearestRestaurantsParams) ([]GetNearestRestaurantsRow, error) {
	rows, err := q.db.Query(ctx, getNearestRestaurants,
		arg.Limit,
		arg.Offset,
		arg.MaxRadius,
		arg.UserID,
		arg.UserLongitude,
		arg.UserLatitude,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetNearestRestaurantsRow{}
	for rows.Next() {
		var i GetNearestRestaurantsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.AddressLine1,
			&i.AddressLine2,
			&i.Postcode,
			&i.Longitude,
			&i.Latitude,
			&i.DistanceMeters,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchRestaurants = `-- name: SearchRestaurants :many
select id,
    name,
    address_line1,
    address_line2,
    postcode,
    ts_rank(
        search_vector,
        websearch_to_tsquery($3::text)
    ) as rank
from restaurants
where search_vector @@ websearch_to_tsquery('english', $3::text)
    and (
        (
            r.valid = true
            and r.created_by is null
        )
        or r.created_by = $4::uuid
    )
order by rank desc
limit $1 offset $2
`

type SearchRestaurantsParams struct {
	Limit      int32     `json:"limit"`
	Offset     int32     `json:"offset"`
	SearchTerm string    `json:"search_term"`
	UserID     uuid.UUID `json:"user_id"`
}

type SearchRestaurantsRow struct {
	ID           uuid.UUID   `json:"id"`
	Name         string      `json:"name"`
	AddressLine1 pgtype.Text `json:"address_line1"`
	AddressLine2 pgtype.Text `json:"address_line2"`
	Postcode     pgtype.Text `json:"postcode"`
	Rank         float32     `json:"rank"`
}

func (q *Queries) SearchRestaurants(ctx context.Context, arg SearchRestaurantsParams) ([]SearchRestaurantsRow, error) {
	rows, err := q.db.Query(ctx, searchRestaurants,
		arg.Limit,
		arg.Offset,
		arg.SearchTerm,
		arg.UserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SearchRestaurantsRow{}
	for rows.Next() {
		var i SearchRestaurantsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.AddressLine1,
			&i.AddressLine2,
			&i.Postcode,
			&i.Rank,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}