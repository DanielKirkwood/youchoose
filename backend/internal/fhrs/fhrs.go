package fhrs

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/DanielKirkwood/youchoose/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EstablishmentDetail represents an individual establishment in the FHRS data.
type EstablishmentDetail struct {
	FHRSID                   int64   `xml:"FHRSID"`
	LocalAuthorityBusinessID string  `xml:"LocalAuthorityBusinessID"`
	BusinessName             string  `xml:"BusinessName"`
	BusinessType             string  `xml:"BusinessType"`
	BusinessTypeID           int     `xml:"BusinessTypeID"`
	AddressLine1             string  `xml:"AddressLine1,omitempty"`
	AddressLine2             string  `xml:"AddressLine2,omitempty"`
	AddressLine3             string  `xml:"AddressLine3,omitempty"`
	PostCode                 string  `xml:"PostCode"`
	RatingValue              string  `xml:"RatingValue"`
	RatingKey                string  `xml:"RatingKey"`
	RatingDate               string  `xml:"RatingDate"`
	LocalAuthorityCode       int     `xml:"LocalAuthorityCode"`
	LocalAuthorityName       string  `xml:"LocalAuthorityName"`
	LocalAuthorityWebsite    string  `xml:"LocalAuthorityWebSite"`
	LocalAuthorityEmail      string  `xml:"LocalAuthorityEmailAddress"`
	SchemeType               string  `xml:"SchemeType"`
	NewRatingPending         bool    `xml:"NewRatingPending"`
	Geocode                  Geocode `xml:"Geocode"`
}

type Geocode struct {
	Longitude float64 `xml:"Longitude"`
	Latitude  float64 `xml:"Latitude"`
}

// FHRSEstablishment represents the root XML structure.
type FHRSEstablishment struct {
	XMLName        xml.Name              `xml:"FHRSEstablishment"`
	Establishments []EstablishmentDetail `xml:"EstablishmentCollection>EstablishmentDetail"`
}

// FetchFHRSData downloads the FHRS XML for the given region.
func FetchFHRSData(apiBaseURL string, regionID string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/%sen-GB.xml", apiBaseURL, regionID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}
	return resp.Body, nil
}

// ParseFHRSData parses the XML data into a structured format.
func ParseFHRSData(data io.Reader) ([]EstablishmentDetail, error) {
	var fhrs FHRSEstablishment
	if err := xml.NewDecoder(data).Decode(&fhrs); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}
	return fhrs.Establishments, nil
}

// FilterRestaurants filters the establishments for only restaurants or cafes.
func FilterRestaurants(establishments []EstablishmentDetail, filterFunc func(EstablishmentDetail) bool) []EstablishmentDetail {
	var filtered []EstablishmentDetail
	for _, e := range establishments {
		if filterFunc(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// StoreRestaurants stores the filtered establishments in the database.
func StoreRestaurants(ctx context.Context, conn *pgxpool.Pool, establishments []EstablishmentDetail) error {
	querier := db.New(conn)

	var params []db.CreateFHRSRawDataParams
	for _, e := range establishments {
		params = append(params, db.CreateFHRSRawDataParams{
			FhrsID:                   e.FHRSID,
			LocalAuthorityBusinessID: pgtype.Text{String: e.LocalAuthorityBusinessID, Valid: e.LocalAuthorityBusinessID != ""},
			BusinessName:             pgtype.Text{String: e.BusinessName, Valid: e.BusinessName != ""},
			BusinessType:             pgtype.Text{String: e.BusinessType, Valid: e.BusinessType != ""},
			BusinessTypeID:           pgtype.Int4{Int32: int32(e.BusinessTypeID), Valid: true},
			AddressLine1:             pgtype.Text{String: e.AddressLine1, Valid: e.AddressLine1 != ""},
			AddressLine2:             pgtype.Text{String: e.AddressLine2, Valid: e.AddressLine2 != ""},
			AddressLine3:             pgtype.Text{String: e.AddressLine3, Valid: e.AddressLine3 != ""},
			Postcode:                 pgtype.Text{String: e.PostCode, Valid: e.PostCode != ""},
			RatingValue:              pgtype.Text{String: e.RatingValue, Valid: e.RatingValue != ""},
			RatingKey:                pgtype.Text{String: e.RatingKey, Valid: e.RatingKey != ""},
			RatingDate:               pgtype.Text{String: e.RatingDate, Valid: e.RatingDate != ""},
			LocalAuthorityCode:       pgtype.Int4{Int32: int32(e.LocalAuthorityCode), Valid: true},
			LocalAuthorityName:       pgtype.Text{String: e.LocalAuthorityName, Valid: e.LocalAuthorityName != ""},
			LocalAuthorityWebsite:    pgtype.Text{String: e.LocalAuthorityWebsite, Valid: e.LocalAuthorityWebsite != ""},
			LocalAuthorityEmail:      pgtype.Text{String: e.LocalAuthorityEmail, Valid: e.LocalAuthorityEmail != ""},
			SchemeType:               pgtype.Text{String: e.SchemeType, Valid: e.SchemeType != ""},
			NewRatingPending:         pgtype.Bool{Bool: e.NewRatingPending, Valid: true},
			Longitude:                pgtype.Float8{Float64: e.Geocode.Longitude, Valid: true},
			Latitude:                 pgtype.Float8{Float64: e.Geocode.Latitude, Valid: true},
		})
	}

	batchResults := querier.CreateFHRSRawData(ctx, params)
	defer batchResults.Close()

	batchResults.Exec(func(i int, err error) {
		if err != nil {
			fmt.Errorf("failed to execute batch at index %d: %w", i, err)
		}
	})

	return nil
}

func SyncRestaurants(ctx context.Context, conn *pgxpool.Pool) error {
	// Fetch all rows from the fhrs_raw_data table
	rows, err := conn.Query(ctx, "select * from fhrs_raw_data where postcode is not null and (longitude <> 0 and latitude <> 0) and (business_type = 'Restaurant/Cafe/Canteen' or business_type = 'Pub/bar/nightclub' or business_type = 'Takeaway/sandwich shop');")
	if err != nil {
		return fmt.Errorf("failed to query fhrs_raw_data: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fhrsData db.FhrsRawDatum
		if err := rows.Scan(&fhrsData.FhrsID, &fhrsData.LocalAuthorityBusinessID, &fhrsData.BusinessName, &fhrsData.BusinessType, &fhrsData.BusinessTypeID, &fhrsData.AddressLine1, &fhrsData.AddressLine2, &fhrsData.AddressLine3, &fhrsData.Postcode, &fhrsData.RatingValue, &fhrsData.RatingKey, &fhrsData.RatingDate, &fhrsData.LocalAuthorityCode, &fhrsData.LocalAuthorityName, &fhrsData.LocalAuthorityWebsite, &fhrsData.LocalAuthorityEmail, &fhrsData.SchemeType, &fhrsData.NewRatingPending, &fhrsData.Longitude, &fhrsData.Latitude, &fhrsData.Created, &fhrsData.Updated); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		// Check if the restaurant exists in the restaurants table
		var existingRestaurant db.Restaurant
		err := conn.QueryRow(ctx, "SELECT id, fhrs_id, name, address_line1, address_line2, address_line3, address_line4, postcode, latitude, longitude, business_type, valid created, updated FROM restaurants WHERE fhrs_id=$1", fhrsData.FhrsID).Scan(
			&existingRestaurant.ID, &existingRestaurant.FhrsID, &existingRestaurant.Name, &existingRestaurant.AddressLine1,
			&existingRestaurant.AddressLine2, &existingRestaurant.AddressLine3, &existingRestaurant.AddressLine4,
			&existingRestaurant.Postcode, &existingRestaurant.Latitude, &existingRestaurant.Longitude,
			&existingRestaurant.BusinessType, &existingRestaurant.Created, &existingRestaurant.Updated)

		if err != nil {
			if err.Error() == "no rows in result set" {
				// If no record exists, insert a new restaurant record
				_, err = conn.Exec(ctx, `
				INSERT INTO restaurants (fhrs_id, name, address_line1, address_line2, address_line3, postcode, latitude, longitude, business_type, valid)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
					fhrsData.FhrsID, fhrsData.BusinessName, fhrsData.AddressLine1, fhrsData.AddressLine2,
					fhrsData.AddressLine3, fhrsData.Postcode, fhrsData.Latitude, fhrsData.Longitude, fhrsData.BusinessType, true)
				if err != nil {
					return fmt.Errorf("failed to insert new restaurant: %v", err)
				}
			} else {
				return fmt.Errorf("failed to check existing restaurant: %v", err)
			}
		} else {
			// If record exists, check if any fields need to be updated
			// You can skip the comparison for custom fields if you don't want to override them
			if existingRestaurant.AddressLine1 != fhrsData.AddressLine1 || existingRestaurant.AddressLine2 != fhrsData.AddressLine2 || existingRestaurant.Postcode != fhrsData.Postcode {
				_, err = conn.Exec(ctx, `
				UPDATE restaurants
				SET name = $1, address_line1 = $2, address_line2 = $3, address_line3 = $4, postcode = $5, latitude = $6, longitude = $7, business_type = $8, updated = now()
				WHERE fhrs_id = $9`,
					fhrsData.BusinessName, fhrsData.AddressLine1, fhrsData.AddressLine2, fhrsData.AddressLine3,
					fhrsData.Postcode, fhrsData.Latitude, fhrsData.Longitude, fhrsData.BusinessType, fhrsData.FhrsID)
				if err != nil {
					return fmt.Errorf("failed to update restaurant: %v", err)
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error encountered during row iteration: %v", err)
	}

	return nil
}
