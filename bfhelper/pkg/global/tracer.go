package global

import (
	"github.com/KomeiDiSanXian/BFHelper/kanban/banner"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer 全局追踪
var Tracer = otel.GetTracerProvider().Tracer(
	"github.com/KomeiDiSanXian/BFHelper",
	trace.WithInstrumentationVersion(banner.Version),
	trace.WithSchemaURL(semconv.SchemaURL),
)
