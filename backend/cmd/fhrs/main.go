package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/DanielKirkwood/youchoose/internal/fhrs"
	"github.com/jackc/pgx/v5"

	"github.com/danielgtaylor/huma/v2/humacli"
)

type Options struct {
	Debug       bool   `doc:"Enable debug logging"`
	FhrsBaseUrl string `doc:"Base URL to fetch FHRS data from" default:"https://ratings.food.gov.uk/api/open-data-files"`
	RegionId    string `doc:"ID of the region to fetch data for" default:"FHRS776"`
	DatabaseURI string `doc:"The database connection string"`
}

func main() {
	cli := humacli.New(func(h humacli.Hooks, opts *Options) {
		// Database connection setup
		ctx := context.Background()
		conn, err := pgx.Connect(ctx, opts.DatabaseURI)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
			os.Exit(1)
		}
		defer conn.Close(ctx)

		err = conn.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
			os.Exit(1)
		}

		// Fetch FHRS data
		body, err := fhrs.FetchFHRSData(opts.FhrsBaseUrl, opts.RegionId)
		if err != nil {
			log.Fatalf("error fetching FHRS data: %v", err)
			os.Exit(1)
		}
		defer body.Close()

		// Parse FHRS data
		establishments, err := fhrs.ParseFHRSData(body)
		if err != nil {
			log.Fatalf("error parsing FHRS data: %v", err)
			os.Exit(1)
		}

		// Filter restaurants
		restaurants := fhrs.FilterRestaurants(establishments, func(e fhrs.EstablishmentDetail) bool {
			return e.BusinessType == "Restaurant/Cafe/Canteen" || e.BusinessType == "Pub/bar/nightclub"
		})

		// Store in database
		if err := fhrs.StoreRestaurants(ctx, conn, restaurants); err != nil {
			log.Fatalf("error storing restaurants: %v", err)
			os.Exit(1)
		}

		fmt.Println("Successfully fetched, filtered, and stored restaurant data.")
		os.Exit(0)
	})

	cli.Run()
}
