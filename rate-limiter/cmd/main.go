package main

import (
	"context"
	"fmt"
	"giovani-milanez/go-expert/rate-limiter/configs"
	"giovani-milanez/go-expert/rate-limiter/internal/limiter"
	"giovani-milanez/go-expert/rate-limiter/internal/middleware"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	storage := limiter.NewRedisRateLimiterStorage(context.Background(), configs.RedisURL)
	err = storage.Connect()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
		return
	}
	_ = storage.Flush()


	// storage          := limiter.NewMemoryRateLimiterStorage()
	ipRateLimiter    := limiter.NewRateLimiter(storage, configs.IpMaxReqPerSecond, configs.IpBlockDuration)
	tokenRateLimiter := limiter.NewRateLimiter(storage, configs.TokenMaxReqPerSecond, configs.TokenBlockDuration)

	http.HandleFunc("/", middleware.RateLimiterMiddleware(tokenRateLimiter, ipRateLimiter, http.HandlerFunc(helloHandler)))
	fmt.Printf("Rate Limiter is running on port %s...\n", configs.ServerPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", configs.ServerPort), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}