package server

import (
	"sync"
	"time"
)

type rateLimiter struct {
	requests map[string]*clientRequests
	mu       sync.RWMutex
}

type clientRequests struct {
    count     int
    firstSeen time.Time
}

func NewRateLimiter() *rateLimiter {
	limiter := &rateLimiter{
		requests: make(map[string]*clientRequests),
	}

	go limiter.cleanupRoutine()

	return limiter
}

func (rl *rateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		rl.cleanup()
	}
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	threshold := time.Now().Add(-1 * time.Hour)

	for ip, client := range rl.requests {
		if client.firstSeen.Before(threshold) {
			delete(rl.requests, ip)
		}
	}
}

func (rl *rateLimiter) Allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    if _, exists := rl.requests[ip]; !exists {
        rl.requests[ip] = &clientRequests{
            count:     0,
            firstSeen: time.Now(),
        }
    }
    
    client := rl.requests[ip]
    
    if time.Since(client.firstSeen) > time.Hour {
        client.count = 0
        client.firstSeen = time.Now()
    }
    
    if client.count >= 100 {
        return false
    }
    
    client.count++
    return true
}

