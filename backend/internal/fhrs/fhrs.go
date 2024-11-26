package fhrs

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// EstablishmentDetail represents an individual establishment in the FHRS data.
type EstablishmentDetail struct {
	FHRSID       string  `xml:"FHRSID"`
	BusinessName string  `xml:"BusinessName"`
	BusinessType string  `xml:"BusinessType"`
	AddressLine1 string  `xml:"AddressLine1"`
	AddressLine2 string  `xml:"AddressLine2"`
	AddressLine3 string  `xml:"AddressLine3"`
	AddressLine4 string  `xml:"AddressLine4"`
	PostCode     string  `xml:"PostCode"`
	Latitude     float64 `xml:"Geocode>Latitude"`
	Longitude    float64 `xml:"Geocode>Longitude"`
}

// FHRSEstablishment represents the root XML structure.
type FHRSEstablishment struct {
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
func StoreRestaurants(ctx context.Context, db *pgx.Conn, establishments []EstablishmentDetail) error {
	query := `
		INSERT INTO restaurants (id, name, address_line1, address_line2, address_line3, postcode, latitude, longitude, business_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE
		SET
			name = EXCLUDED.name,
			address_line1 = EXCLUDED.address_line1,
			address_line2 = EXCLUDED.address_line2,
			address_line3 = EXCLUDED.address_line3,
			postcode = EXCLUDED.postcode,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			business_type = EXCLUDED.business_type;
	`

	batch := &pgx.Batch{}
	for _, est := range establishments {
		batch.Queue(query, est.FHRSID, est.BusinessName, est.AddressLine1, est.AddressLine2, est.AddressLine3, est.PostCode, est.Latitude, est.Longitude, est.BusinessType)
	}

	batchResults := db.SendBatch(ctx, batch)
	defer batchResults.Close()

	for i := 0; i < len(establishments); i++ {
		_, err := batchResults.Exec() // Exec checks for errors on each query
		if err != nil {
			return fmt.Errorf("batch query %d failed: %w", i, err)
		}
	}

	return nil
}
