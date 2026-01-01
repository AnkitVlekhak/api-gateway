package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnkitVlekhak/api-gateway/internal/config"
	"github.com/AnkitVlekhak/api-gateway/internal/gateway"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	gateway, err := gateway.NewGateway(config)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{Addr: ":8080", Handler: gateway}

	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("🚀 API Gateway running on :8080")
		if err := server.Serve(listner); err != nil {
			log.Fatal(err)
		}
	}()

	// 5. Graceful shutdown handling
	shutdownCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-shutdownCtx.Done()
	log.Println("🛑 Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	} else {
		log.Println("✅ Server shut down gracefully")
	}

}
