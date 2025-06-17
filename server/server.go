package server

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	httpServer  *http.Server
	db          *sql.DB
	middlewares []func(http.Handler) http.Handler
}

func (this *Server) UseMiddleware(middleware ...func(http.Handler) http.Handler) {
	this.middlewares = append(this.middlewares, middleware...)
}

func (this *Server) Run(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	//setup the handler
	mux := http.NewServeMux()
	handler := NewLibraryHandler(this.db)
	handler.RegisterRoutes(mux)

	//Middlewares
	this.httpServer.Handler = recoverMiddleware(corsMiddleware(mux))

	for _, m := range this.middlewares {
		this.httpServer.Handler = m(this.httpServer.Handler)
	}

	//Signal interrupt
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer func() {
		stop()
		log.Printf("Server.httpServer.Shutdown(context.Background()): %v\n", this.httpServer.Shutdown(context.Background()))
	}()
	this.httpServer.BaseContext = func(_ net.Listener) context.Context { return ctx }

	serverErrChan := make(chan error)
	go func() {
		log.Println("Server started!")
		serverErrChan <- this.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("Interrupt signal recieved!")
		return
	case err := <-serverErrChan:
		log.Print("Server error:", err)
		return
	}
}

func NewServer(port string, db *sql.DB) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			ReadTimeout:  time.Second,
			WriteTimeout: 10 * time.Second,
		},
		db:          db,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}
