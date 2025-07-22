package limiter

import (
	// "sync"
	"time"
)

type MemoryRateLimiterStorage struct {
	// sync.RWMutex
	rates map[string]*Rate
}

func NewMemoryRateLimiterStorage() *MemoryRateLimiterStorage {
	return &MemoryRateLimiterStorage{
		rates: make(map[string]*Rate),
	}
}

func (m *MemoryRateLimiterStorage) Flush() error {
	m.rates = make(map[string]*Rate)
	return nil
}

func (m *MemoryRateLimiterStorage) Get(key string) (*Rate, error) {
	// m.Lock()
	// defer m.Unlock()
	rate, exists := m.rates[key]
	if !exists {
		rate = &Rate{Count: 1, FirstSeen: time.Now(), BlockedUntil: time.Now().Add(-1 * time.Hour)}
		m.rates[key] = rate
	}
	return rate, nil
}

func (m *MemoryRateLimiterStorage) Update(key string, rate *Rate) error {
	// m.Lock()
	// defer m.Unlock()
	m.rates[key] = rate
	return nil
}