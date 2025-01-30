package main

import (
	"os"

	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/config"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/internal/adapters/db"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/internal/adapters/grpc"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/internal/adapters/payment"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/internal/application/core/api"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	service     = "order"
	environment = "dev"
	id          = 1
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp), // export the tracing matrics in batch
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,                      // the url that contains the opentelemetry schema
			semconv.ServiceNameKey.String(service), // an attribute that describes the server name
			attribute.String("environment", environment),
			attribute.Int64("ID", id), // an arbitrary id that can be used in tracing the dashboard
		)),
	)
	return tp, nil
}

func init() {
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{FieldMap: log.FieldMap{"msg": "message"}},
		// rename the log fields
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type customLogger struct {
	formatter log.JSONFormatter
}

func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	//Below injection is Just to understand what Context has
	entry.Data["Context"] = span.SpanContext()
	return l.formatter.Format(entry) // Injects the trace and span into the current log message
}

func main() {
	tp, err := tracerProvider("http://jaeger-otel.jaeger.svc.cluster.local:14278/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp) // set the tracing provider through the opentelemetry sdk
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	// In the tracing provider configuration, we simply provide a tracing collect endpoint,
	// which we installed via the Helm Chart, and set the tracing provider using OpenTelem-
	// etry SDK. Then we configure the propagation strategy to propagate traces and spans
	// from one service to another. For example, since the Order service calls the Payment
	// service to charge a customer, existing trace metadata will be propagated to the Payment
	// service to see the whole request flow in the Jaeger tracing UI

	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize payment stub. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
