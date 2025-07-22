package limiter

import (
	// "log"
	// "fmt"
	"math"
	"time"
)

type Rate struct {
	Count        uint16    `json:"count"`
	FirstSeen    time.Time `json:"first_seen"`
	BlockedUntil time.Time `json:"blocked_until"`
	FlagReset bool   `json:"flag_reset"`
}

func (r *Rate) GetReqPerSecond() uint {
	// round up to the next second
	secondsSinceFirstSeen := uint(math.Trunc(time.Since(r.FirstSeen).Seconds())) + 1
	// fmt.Println(fmt.Sprintf("Seconds since first seen: %d", secondsSinceFirstSeen))
	reqpers := uint(r.Count) / secondsSinceFirstSeen
	// fmt.Println(fmt.Sprintf("Requests per second: %s - %d", reqpers))
	return reqpers
}

func (r *Rate) Increment() {
	r.Count++
}

func (r *Rate) Block(duration time.Duration) {
	r.BlockedUntil = time.Now().Add(duration)
}

func (r *Rate) IsBlocked() bool {
	return time.Now().Before(r.BlockedUntil)
}

func (r *Rate) NeedsReset() bool {
	return r.FlagReset
}
