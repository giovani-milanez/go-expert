package limiter

import (
	// "fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	sync.RWMutex
	storage              RateLimiterStorage
	MaxRequestsPerSecond uint
	BlockTime            time.Duration
	rates                map[string]*Rate
}

func NewRateLimiter(storage RateLimiterStorage, maxRequestsPerSecond uint, blockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		storage:              storage,
		MaxRequestsPerSecond: maxRequestsPerSecond,
		BlockTime:            blockTime,
		rates:                make(map[string]*Rate),
	}
}

func (m *RateLimiter) IsAllowed(key string) (bool, error) {
	m.Lock()
	defer m.Unlock()

	rate, err := m.storage.Get(key)
	if err != nil {
		return false, err
	}
	
	if rate.IsBlocked() {
		if !rate.FlagReset {
			rate.FlagReset = true
			err = m.storage.Update(key, rate)
			if err != nil {
				return false, err
			}
		}
		return false, nil
	}
	reqps := rate.GetReqPerSecond()
	// fmt.Println(fmt.Sprintf("Requests per second for key %s: %d", key, reqps))
	if reqps > m.MaxRequestsPerSecond {
		// fmt.Println("Blocking key:", key, "for exceeding max requests per second")
		rate.Block(m.BlockTime)
		err = m.storage.Update(key, rate)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	if rate.NeedsReset() {
		// fmt.Println("Resetting rate for key:", key)
		rate.Count = 1
		rate.FirstSeen = time.Now()
		rate.BlockedUntil = time.Now().Add(-1 * time.Hour)
		rate.FlagReset = false	
	}
	
	rate.Increment()	
	err = m.storage.Update(key, rate)
	if err != nil {
		return false, err
	}
	return true, nil
}