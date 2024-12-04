package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DanielKirkwood/youchoose/internal/api"
	"github.com/DanielKirkwood/youchoose/internal/db"
	"github.com/DanielKirkwood/youchoose/internal/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/jackc/pgx/v5"
)

type Options struct {
	Port        int    `doc:"Port to listen on" short:"p" default:"8888"`
	Version     string `doc:"The api version" short:"v" default:"1.0.0"`
	DatabaseURI string `doc:"The database connection string"`
}

type HelloOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		ctx := context.Background()

		router := http.NewServeMux()
		humaApi := humago.New(router, huma.DefaultConfig("You Choose API", opts.Version))

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", opts.Port),
			Handler: router,
		}

		conn, err := pgx.Connect(ctx, opts.DatabaseURI)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
			os.Exit(1)
		}

		err = conn.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
			os.Exit(1)
		}

		queries := db.New(conn)

		api.RegisterRestaurantRoutes(humaApi, &services.RestaurantService{Queries: queries})

		hooks.OnStart(func() {
			fmt.Printf("Starting server version %s on port %s ...\n", opts.Version, server.Addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("server error: %v", err)
			}
		})

		hooks.OnStop(func() {
			// Gracefully shutdown your server here
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Fatalf("server shutdown error: %v", err)
			}

			if err := conn.Close(ctx); err != nil {
				log.Fatalf("failed to close database connection: %v", err)
			}
		})
	})

	cli.Run()
}
