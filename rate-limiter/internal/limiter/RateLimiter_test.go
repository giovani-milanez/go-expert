package limiter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterNoBlock(t *testing.T) {
	MAX_REQUESTS_PER_SECOND := uint(10)
	BLOCK_TIME := 10 * time.Second
	REQUEST_COUNT := 10
	storage := NewMemoryRateLimiterStorage()
	limiter := NewRateLimiter(storage, MAX_REQUESTS_PER_SECOND, BLOCK_TIME)
	for i := 0; i < REQUEST_COUNT; i++ {
		if ok,err := limiter.IsAllowed("user1"); err != nil || !ok {
			t.Errorf("Request %d was blocked: %v", i, err)
		}
	}
}

func TestRateLimiterBlock(t *testing.T) {
	MAX_REQUESTS_PER_SECOND := uint(10)
	BLOCK_TIME := 10 * time.Second
	REQUEST_COUNT := 10
	storage := NewMemoryRateLimiterStorage()
	limiter := NewRateLimiter(storage, MAX_REQUESTS_PER_SECOND, BLOCK_TIME)
	for i := 0; i < REQUEST_COUNT; i++ {
		if ok,err := limiter.IsAllowed("user1"); err != nil || !ok {
			t.Errorf("Request %d was blocked: %v", i, err)
		}
	}
	ok, err := limiter.IsAllowed("user1");
	if err != nil {
		t.Errorf("Error checking rate limit after %d requests: %v", REQUEST_COUNT, err)
	}
	if ok {
		t.Errorf("Request is not block but it should be")
	}
}

func TestRateLimiterMultipleRequests(t *testing.T) {
	MAX_REQUESTS_PER_SECOND := uint(10)
	BLOCK_TIME := 1 * time.Second
	REQUEST_COUNT := 100
	WAIT_TIME := 50 * time.Millisecond
	storage := NewMemoryRateLimiterStorage()
	limiter := NewRateLimiter(storage, MAX_REQUESTS_PER_SECOND, BLOCK_TIME)

	allowedCount := 0
	rejectedCount := 0

	for i := 0; i < REQUEST_COUNT; i++ {
		ok, err := limiter.IsAllowed("user1")
		if err != nil {
			t.Errorf("Error checking rate limit: %v\n", err)
		}
		if ok {
			allowedCount++
		} else {
			rejectedCount++
		}
		time.Sleep(WAIT_TIME)
	}
	if allowedCount != 40 {
		t.Errorf("Expected 40 allowed requests, got %d", allowedCount)
	}
	if rejectedCount != 60 {
		t.Errorf("Expected 60 rejected requests, got %d", rejectedCount)
	}
}

func TestRateLimiterMultipleRequestsThreadSameKey(t *testing.T) {
	MAX_REQUESTS_PER_SECOND := uint(10)
	BLOCK_TIME := 1 * time.Second
	REQUEST_COUNT := 100
	WAIT_TIME := 50 * time.Millisecond
	storage := NewMemoryRateLimiterStorage()
	limiter := NewRateLimiter(storage, MAX_REQUESTS_PER_SECOND, BLOCK_TIME)

	allowedCount := 0
	rejectedCount := 0
	
	var wg sync.WaitGroup
	var mu sync.RWMutex
	wg.Add(2)

	workerFn := func(key string, n int) {
		defer wg.Done()
		for i := 0; i < n; i++ {
			ok, err := limiter.IsAllowed(key)
			if err != nil {
				t.Errorf("Error checking rate limit: %v\n", err)
			}
			mu.Lock()
			if ok {
				allowedCount++
			} else {
				rejectedCount++
			}
			mu.Unlock()
			time.Sleep(WAIT_TIME)
		}
	}

	go workerFn("user1", REQUEST_COUNT)
	go workerFn("user1", REQUEST_COUNT)

	wg.Wait()

	if allowedCount != 40 {
		t.Errorf("Expected 40 allowed requests, got %d", allowedCount)
	}
	if rejectedCount != 160 {
		t.Errorf("Expected 160 rejected requests, got %d", rejectedCount)
	}
}

func TestRateLimiterMultipleRequestsThreadDifferentKeys(t *testing.T) {
	MAX_REQUESTS_PER_SECOND := uint(10)
	BLOCK_TIME := 1 * time.Second
	REQUEST_COUNT := 100
	WAIT_TIME := 50 * time.Millisecond
	storage := NewMemoryRateLimiterStorage()
	limiter := NewRateLimiter(storage, MAX_REQUESTS_PER_SECOND, BLOCK_TIME)

	allowedCount := 0
	rejectedCount := 0
	var mu sync.RWMutex

	threads := 100
	
	var wg sync.WaitGroup
	wg.Add(threads)

	workerFn := func(key string, n int) {
		defer wg.Done()
		for i := 0; i < n; i++ {
			ok, err := limiter.IsAllowed(key)
			if err != nil {
				t.Errorf("Error checking rate limit: %v\n", err)
			}
			mu.Lock()
			if ok {
				allowedCount++
			} else {
				rejectedCount++
			}
			mu.Unlock()
			time.Sleep(WAIT_TIME)
		}
	}

	for i := 0; i < threads; i++ {
		go workerFn(fmt.Sprintf("user%d", i+1), REQUEST_COUNT)
	}

	wg.Wait()

	if allowedCount != 40 * threads {
		t.Errorf("Expected %d allowed requests, got %d", 40 * threads, allowedCount)
	}
	if rejectedCount != 60 * threads {
		t.Errorf("Expected %d rejected requests, got %d", 60 * threads, rejectedCount)
	}
}