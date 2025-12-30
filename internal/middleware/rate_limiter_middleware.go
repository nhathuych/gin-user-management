package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mut     sync.Mutex
	clients = make(map[string]*Client)
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mut.Lock()
	defer mut.Unlock()

	client, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(3, 10) // 3 requests/second with a burst of up to 10
		newClient := &Client{limiter, time.Now()}
		clients[ip] = newClient

		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mut.Lock()

		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mut.Unlock()
	}
}

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)
		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		ctx.Next()
	}
}
