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
	"github.com/AnkitVlekhak/api-gateway/internal/gateway/builder"
	"github.com/redis/go-redis/v9"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	routes, err := builder.BuildRoutes(config)
	if err != nil {
		log.Fatal(err)
	}

	router := gateway.NewPrefixRouter(routes)

	trustedProxyIdentityResolver, err := gateway.NewTrustedProxyIdentityResolver([]string{
		"10.0.0.0/8",
		"192.168.0.0/16",
	})
	if err != nil {
		log.Fatal(err)
	}

	apiKeyIdentityResolver := gateway.NewAPIKeyIdentityResolver("X-API-Key")

	composite := gateway.NewCompositeIdentityResolver(
		apiKeyIdentityResolver,        // highest priority
		trustedProxyIdentityResolver,  // trusted proxy IP
		&gateway.IPIdentityResolver{}, // fallback
	)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ratelimiter := gateway.NewRedisTokenBucketRateLimiter(redisClient)

	gateway, err := gateway.NewGateway(router, composite, ratelimiter)
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
