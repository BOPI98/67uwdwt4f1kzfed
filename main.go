package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	server *http.Server
	db     *sql.DB
)

func init() {
	initConfig()

	var err error
	db, err = sql.Open("mysql", config.DbUser+":"+config.DbPassword+"@tcp("+config.DbHost+")/"+config.DbName+"?parseTime=true")
	if err != nil {
		log.Fatal("Connecting to the database failed:", err)
	}

	server = &http.Server{
		Addr:         ":" + config.ApiPort,
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return context.Background()
		},
	}

	initOtelProviders()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	defer func() {
		tracerProvider.Shutdown(ctx)
		metricProvider.Shutdown(ctx)
		db.Close()
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/books", Handle_Books)
	mux.HandleFunc("/books/{id}/borrow", Handle_Books_Id_Borrow)

	mux.HandleFunc("/borrowers", Handle_Borrowers)
	mux.HandleFunc("/borrowers/{id}", Handle_Borrowers_Id)
	mux.HandleFunc("/borrowers/{id}/borrowedbooks", Handle_Borrowers_Id_BorrowedBooks)
	server.Handler = recoverMiddleware(corsMiddleware(mux))
	server.BaseContext = func(_ net.Listener) context.Context { return ctx }

	serverErrChan := make(chan error)
	go func() {
		serverErrChan <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		stop()
	case err := <-serverErrChan:
		log.Print(err)
		return
	}
}
