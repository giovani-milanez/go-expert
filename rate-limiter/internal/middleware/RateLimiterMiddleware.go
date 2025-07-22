package middleware

import (
	"giovani-milanez/go-expert/rate-limiter/internal/limiter"
	"net"
	"net/http"
)

func RateLimiterMiddleware(tokenRl *limiter.RateLimiter, ipRl *limiter.RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Primeiro vefifica o Token, que tem prioridade sobre o IP		
		apiKey := r.Header.Get("API_KEY")
		if apiKey != "" {
			ok, err := tokenRl.IsAllowed(apiKey)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !ok {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Se o Token não existir, verifica o IP
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Error getting IP address", http.StatusInternalServerError)
			return
		}
		ok, err := ipRl.IsAllowed(host)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}