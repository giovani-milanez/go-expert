package limiter

type RateLimiterStorage interface {	
	Get(key string) (*Rate, error)
	Update(key string, rate *Rate) error
	// clear all data
	Flush() error
}