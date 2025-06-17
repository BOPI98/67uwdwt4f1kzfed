package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/BOPI98/67uwdwt4f1kzfed/server"
	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const SERVER_NAME = "library-management"

var (
	tracer trace.Tracer
	meter  metric.Meter
	logger *slog.Logger

	memoryHeapGauge   metric.Int64ObservableGauge
	durationHistogram metric.Float64Histogram
)

func init() {
	//setup otel
	initOtelProviders()

	//read config
	initConfig()

	var err error
	memoryHeapGauge, err = meter.Int64ObservableGauge("memory.heap",
		metric.WithDescription("Allocated heap objects"),
		metric.WithInt64Callback(func(ctx context.Context, io metric.Int64Observer) error {
			m := runtime.MemStats{}
			runtime.ReadMemStats(&m)
			io.Observe(int64(m.HeapAlloc))
			return nil
		}))
	if err != nil {
		log.Fatal(err)
	}

	durationHistogram, err = meter.Float64Histogram("request.duration",
		metric.WithDescription("The duration of request execution"),
		metric.WithUnit("s"),
	)
}

func main() {
	defer func() {
		log.Printf("tracerProvider.Shutdown(context.Background()): %v\n", tracerProvider.Shutdown(context.Background()))
		log.Printf("metricProvider.Shutdown(context.Background()): %v\n", metricProvider.Shutdown(context.Background()))
		log.Printf("logProvider.Shutdown(context.Background()): %v\n", logProvider.Shutdown(context.Background()))
	}()

	db, err := otelsql.Open("mysql", config.DbUser+":"+config.DbPassword+"@tcp("+config.DbHost+")/"+config.DbName+"?parseTime=true")
	if err != nil {
		log.Fatal("Connecting to the database failed:", err)
	}
	defer db.Close()
	log.Println("Database connected!")

	srv := server.NewServer(config.ApiPort, db)
	srv.UseMiddleware(otelMiddleware, metricMiddleware)
	srv.Run(nil)
}

func otelMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewMiddleware("http.handler")(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.Method + " " + r.URL.String())
			next.ServeHTTP(w, r)
		},
	))
}

func metricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		durationHistogram.Record(context.Background(), float64(time.Since(start).Seconds()))
	})
}
