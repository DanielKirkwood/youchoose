package server

import "net/http"

func addRoutes(mux *http.ServeMux) {
	mux.Handle("/hello", handleHello())
}

func handleHello() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello World!"))
		},
	)
}
