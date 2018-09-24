package server

import (
	"context"
	"fmt"
	"net/http"
)

type HTTPServer struct {
	Server *http.Server
}

// Start initialize the HTTP server
func Start(handler http.Handler) *HTTPServer {

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		WriteTimeout: config.httpWriteTimeout,
		ReadTimeout:  config.httpReadTimeout,
		IdleTimeout:  config.httpIdleTimeout,
		Handler:      handler,
	}

	go func() {
		fmt.Printf("Starting HTTP Server on port %d\n", config.port)
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error on start server: %v\n", err)
		}
	}()

	return &HTTPServer{
		Server: server,
	}
}

func (h *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), config.httpGracefulShutdown)
	defer cancel()

	h.Server.Shutdown(ctx)
	fmt.Println("HTTP Server shutting down!")
}
