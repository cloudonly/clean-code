package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"

	"github.com/luomu/clean-code/internal/router"
)

func main() {
	// 初始化 tracer
	tp, err := initTracer("http://local.jaeger-collector.com/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	// 确保在程序结束时关闭 tracer provider
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	addr := ":8080"
	router := router.NewGinEngine()
	// use gin.Engine.Handler() to support h2c
	Serve(addr, router.Handler())
}

func Serve(addr string, handler http.Handler) {
	srv := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

// 初始化 Tracer， 设置采样器，指定资源属性并创建 Jaeger exporter
func initTracer(url string) (*tracesdk.TracerProvider, error) {
	// 创建 jaeger exporter，并指定 endpoint
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	// 创建 tracer provider
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// 采样策略设置为 AlwaysSample，即记录所有 Span
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// 设置资源属性，例如服务名，环境等信息
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("gin-trace-demo"),
			attribute.String("env", "dev"),
		)),
	)
	// 将 tp 设置为全局的 tracer provider
	otel.SetTracerProvider(tp)
	// 设置默认的 TextMapPropagator
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}
