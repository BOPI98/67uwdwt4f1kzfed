package main

import (
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	tracerProvider *sdkTrace.TracerProvider
	metricProvider *sdkMetric.MeterProvider
)

func initOtelProviders() {
	date := time.Now().Format("2006-01-02_15-04-05")
	traceFile, err := os.Create("./traces_" + date + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	traceExporter, err := stdouttrace.New(stdouttrace.WithWriter(traceFile))
	if err != nil {
		log.Fatal(err)
	}

	tracerProvider = sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(traceExporter),
	)

	metricFile, err := os.Create("./metrics_" + date + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	metricExporter, err := stdoutmetric.New(stdoutmetric.WithWriter(metricFile))
	if err != nil {
		log.Fatal(err)
	}

	metricProvider = sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(
			sdkMetric.NewPeriodicReader(metricExporter),
		),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(metricProvider)
}
