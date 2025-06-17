package main

import (
	"log"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	tracerProvider *sdkTrace.TracerProvider
	metricProvider *sdkMetric.MeterProvider
	logProvider    *sdkLog.LoggerProvider
)

func initOtelProviders() {
	date := time.Now().Format("2006-01-02_15-04-05")
	//Tracer
	traceFile, err := os.Create("./traces_" + date + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	traceExporter, err := stdouttrace.New(
		stdouttrace.WithWriter(traceFile),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatal(err)
	}

	tracerProvider = sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(traceExporter),
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
	)

	//Metric
	metricFile, err := os.Create("./metrics_" + date + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	metricExporter, err := stdoutmetric.New(
		stdoutmetric.WithWriter(metricFile),
		stdoutmetric.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatal(err)
	}

	metricProvider = sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(
			sdkMetric.NewPeriodicReader(metricExporter),
		),
	)
	//Logger
	logFile, err := os.Create("./log_" + date + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	logExporter, err := stdoutlog.New(
		stdoutlog.WithWriter(logFile),
		stdoutlog.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatal(err)
	}

	logProvider = sdkLog.NewLoggerProvider(
		sdkLog.WithProcessor(
			sdkLog.NewBatchProcessor(logExporter),
		),
	)

	//Set providers
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(metricProvider)
	global.SetLoggerProvider(logProvider)

	tracer = otel.Tracer(SERVER_NAME)
	meter = otel.Meter(SERVER_NAME)
	logger = otelslog.NewLogger(SERVER_NAME)
}
