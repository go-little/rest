package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HTTPServer struct {
	*http.Server
}

// Start initialize the HTTP server
func Start(config *http.Server) *HTTPServer {

	go func() {
		fmt.Printf("Starting HTTP Server on port %s\n", config.Addr)
		if err := config.ListenAndServe(); err != nil {
			fmt.Printf("Error on start server: %v\n", err)
		}
	}()

	return &HTTPServer{
		Server: config,
	}
}

func (h *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.Shutdown(ctx)
	fmt.Println("HTTP Server shutting down!")
}
