// Package tracer 实现链路追踪
package tracer

import (
	"context"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/kanban/banner"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func newExporter(ctx context.Context, url string) (*otlptrace.Exporter, error) {
	opt := []otlptracehttp.Option{otlptracehttp.WithEndpoint(url)}
	if !global.TraceSetting.UseHTTPS {
		opt = append(opt, otlptracehttp.WithInsecure())
	}
	client := otlptracehttp.NewClient(opt...)
	return otlptrace.New(ctx, client)
}

func newTraceProvider(expo sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("BFHelperService"),
			semconv.ServiceVersion(banner.Version),
		),
	)
	if err != nil {
		return nil, err
	}
	// 总是开启追踪
	sampler := sdktrace.WithSampler(sdktrace.AlwaysSample())
	if !global.TraceSetting.Enabled {
		sampler = sdktrace.WithSampler(sdktrace.NeverSample())
	}
	return sdktrace.NewTracerProvider(sdktrace.WithBatcher(expo), sdktrace.WithResource(r), sampler), nil
}

// InstallExportPipeline 追踪导出
func InstallExportPipeline(ctx context.Context, url string) (func(ctx context.Context) error, error) {
	expo, err := newExporter(ctx, url)
	if err != nil {
		return nil, err
	}
	tp, err := newTraceProvider(expo)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

// AddEventWithDescription 在span添加k-v形式的描述
func AddEventWithDescription(desc ...attribute.KeyValue) trace.EventOption {
	return trace.WithAttributes(desc...)
}

// Description 返回attribute.KeyValue
func Description(k, v string) attribute.KeyValue {
	return attribute.String(k, v)
}
