// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type FhrsRawDatum struct {
	FhrsID                   int64              `json:"fhrs_id"`
	LocalAuthorityBusinessID pgtype.Text        `json:"local_authority_business_id"`
	BusinessName             pgtype.Text        `json:"business_name"`
	BusinessType             pgtype.Text        `json:"business_type"`
	BusinessTypeID           pgtype.Int4        `json:"business_type_id"`
	AddressLine1             pgtype.Text        `json:"address_line1"`
	AddressLine2             pgtype.Text        `json:"address_line2"`
	AddressLine3             pgtype.Text        `json:"address_line3"`
	Postcode                 pgtype.Text        `json:"postcode"`
	RatingValue              pgtype.Text        `json:"rating_value"`
	RatingKey                pgtype.Text        `json:"rating_key"`
	RatingDate               pgtype.Text        `json:"rating_date"`
	LocalAuthorityCode       pgtype.Int4        `json:"local_authority_code"`
	LocalAuthorityName       pgtype.Text        `json:"local_authority_name"`
	LocalAuthorityWebsite    pgtype.Text        `json:"local_authority_website"`
	LocalAuthorityEmail      pgtype.Text        `json:"local_authority_email"`
	SchemeType               pgtype.Text        `json:"scheme_type"`
	NewRatingPending         pgtype.Bool        `json:"new_rating_pending"`
	Longitude                pgtype.Float8      `json:"longitude"`
	Latitude                 pgtype.Float8      `json:"latitude"`
	Created                  pgtype.Timestamptz `json:"created"`
	Updated                  pgtype.Timestamptz `json:"updated"`
}

type Restaurant struct {
	ID           uuid.UUID          `json:"id"`
	FhrsID       int32              `json:"fhrs_id"`
	Name         string             `json:"name"`
	AddressLine1 pgtype.Text        `json:"address_line1"`
	AddressLine2 pgtype.Text        `json:"address_line2"`
	AddressLine3 pgtype.Text        `json:"address_line3"`
	AddressLine4 pgtype.Text        `json:"address_line4"`
	Postcode     pgtype.Text        `json:"postcode"`
	Latitude     pgtype.Float8      `json:"latitude"`
	Longitude    pgtype.Float8      `json:"longitude"`
	BusinessType pgtype.Text        `json:"business_type"`
	Valid        pgtype.Bool        `json:"valid"`
	Created      pgtype.Timestamptz `json:"created"`
	Updated      pgtype.Timestamptz `json:"updated"`
	SearchVector interface{}        `json:"search_vector"`
	Geolocation  interface{}        `json:"geolocation"`
	CreatedBy    pgtype.UUID        `json:"created_by"`
}

type User struct {
	ID          uuid.UUID          `json:"id"`
	DisplayName string             `json:"display_name"`
	Email       pgtype.Text        `json:"email"`
	Created     pgtype.Timestamptz `json:"created"`
	Updated     pgtype.Timestamptz `json:"updated"`
}
