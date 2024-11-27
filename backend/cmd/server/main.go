package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
)

type Options struct {
	Port    int    `doc:"Port to listen on" short:"p" default:"8888"`
	Version string `doc:"The api version" short:"v" default:"1.0.0"`
}

type HelloOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		router := http.NewServeMux()
		api := humago.New(router, huma.DefaultConfig("You Choose API", opts.Version))

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", opts.Port),
			Handler: router,
		}

		huma.Register(api, huma.Operation{
			OperationID: "get-hello",
			Method:      http.MethodGet,
			Path:        "/hello",
			Summary:     "Get a hello",
			Description: "Get a hello message",
			Tags:        []string{"Hello"},
		}, func(ctx context.Context, input *struct{}) (*HelloOutput, error) {
			resp := &HelloOutput{}
			resp.Body.Message = "Hi"
			return resp, nil
		})

		hooks.OnStart(func() {
			fmt.Printf("Starting server version %s on port %s ...\n", opts.Version, server.Addr)
			server.ListenAndServe()
		})

		hooks.OnStop(func() {
			// Gracefully shutdown your server here
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	})

	cli.Run()
}
