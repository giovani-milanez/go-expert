package middleware

import (
	"fmt"
	"giovani-milanez/go-expert/rate-limiter/internal/limiter"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
	// "context"
)

	
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func getStorage(t *testing.T) limiter.RateLimiterStorage {
	// Atencao, tests com Redis vao falhar por causa do delay entre chamadas
	// fica dificil de ter um numero exato de requests permitidas/bloqueadas

	// storage := limiter.NewRedisRateLimiterStorage(context.Background(), "127.0.0.1:6379")
	// err := storage.Connect()
	// if err != nil {
	// 	t.Fatal("Could not connect to Redis:", err)
	// 	return nil
	// }
	// return storage
	return limiter.NewMemoryRateLimiterStorage()
}

func TestMiddlewareIpRateLimitOK(t *testing.T) {
	storage := getStorage(t)
	defer storage.Flush()
	ipRateLimiter := limiter.NewRateLimiter(storage, 10, 60*time.Second)
	tokenRateLimiter := limiter.NewRateLimiter(storage, 100, 30*time.Second)

	handler := RateLimiterMiddleware(tokenRateLimiter, ipRateLimiter, http.HandlerFunc(helloHandler))

	for range 10 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	}
}

func TestMiddlewareIpRateLimitBlock(t *testing.T) {
	storage := getStorage(t)
	defer storage.Flush()
	ipRateLimiter := limiter.NewRateLimiter(storage, 10, 60*time.Second)
	tokenRateLimiter := limiter.NewRateLimiter(storage, 100, 30*time.Second)

	handler := RateLimiterMiddleware(tokenRateLimiter, ipRateLimiter, http.HandlerFunc(helloHandler))

	for range 10 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", w.Code)
	}
}

func TestMiddlewareIpRateLimitMultiple(t *testing.T) {
	storage := getStorage(t)
	defer storage.Flush()
	ipRateLimiter := limiter.NewRateLimiter(storage, 10, 1*time.Second)
	tokenRateLimiter := limiter.NewRateLimiter(storage, 100, 30*time.Second)

	handler := RateLimiterMiddleware(tokenRateLimiter, ipRateLimiter, http.HandlerFunc(helloHandler))
	
	allowedCount := 0
	rejectedCount := 0

	for range 100 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		switch w.Code {
		case http.StatusOK:
			allowedCount++
		case http.StatusTooManyRequests:
			rejectedCount++
		default:
			t.Errorf("Unexpected status code: %v", w.Code)
		}
		time.Sleep(50 * time.Millisecond)
	}

	if allowedCount != 40 {
		t.Errorf("Expected 40 allowed requests, got %d", allowedCount)
	}
	if rejectedCount != 60 {
		t.Errorf("Expected 60 rejected requests, got %d", rejectedCount)
	}
}


func TestMiddlewareTokenRateLimitMultiple(t *testing.T) {
	storage := getStorage(t)
	defer storage.Flush()
	ipRateLimiter := limiter.NewRateLimiter(storage, 10, 1*time.Second)
	tokenRateLimiter := limiter.NewRateLimiter(storage, 15, 2*time.Second)

	handler := RateLimiterMiddleware(tokenRateLimiter, ipRateLimiter, http.HandlerFunc(helloHandler))
	
	allowedCount := 0
	rejectedCount := 0

	for range 100 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("API_KEY", "test-key")
		w := httptest.NewRecorder()

		handler(w, req)

		switch w.Code {
		case http.StatusOK:
			allowedCount++
		case http.StatusTooManyRequests:
			rejectedCount++
		default:
			t.Errorf("Unexpected status code: %v", w.Code)
		}
		time.Sleep(50 * time.Millisecond)
	}

	if allowedCount != 30 {
		t.Errorf("Expected 30 allowed requests, got %d", allowedCount)
	}
	if rejectedCount != 70 {
		t.Errorf("Expected 70 rejected requests, got %d", rejectedCount)
	}
}

func TestMiddlewareTokenRateLimitMultipleThreads(t *testing.T) {
	storage := getStorage(t)
	defer storage.Flush()
	ipRateLimiter := limiter.NewRateLimiter(storage, 5, 1*time.Second)
	tokenRateLimiter := limiter.NewRateLimiter(storage, 10, 1*time.Second)

	handler := RateLimiterMiddleware(tokenRateLimiter, ipRateLimiter, http.HandlerFunc(helloHandler))
	
	allowedCount := 0
	rejectedCount := 0

	var mu sync.RWMutex

	threads := 10
	
	var wg sync.WaitGroup
	wg.Add(threads)

	workerFn := func(apiKey string, n int) {
		defer wg.Done()
		for range n {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add("API_KEY", apiKey)
			w := httptest.NewRecorder()
	
			handler(w, req)
			
			mu.Lock()
			switch w.Code {
			case http.StatusOK:
				allowedCount++
			case http.StatusTooManyRequests:
				rejectedCount++
			default:
				t.Errorf("Unexpected status code: %v", w.Code)
			}
			mu.Unlock()
			time.Sleep(50 * time.Millisecond)
		}
	}

	for i := 0; i < threads; i++ {
		go workerFn(fmt.Sprintf("key%d", i+1), 100)
	}

	wg.Wait()

	if allowedCount != 40 * threads {
		t.Errorf("Expected %d allowed requests, got %d", 40 * threads, allowedCount)
	}
	if rejectedCount != 60 * threads {
		t.Errorf("Expected %d rejected requests, got %d", 60 * threads, rejectedCount)
	}
}