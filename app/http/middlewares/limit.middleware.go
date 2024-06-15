package middlewares

import (
	"context"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

type LimitMiddleware struct{}

// Handler 用来做限流
func (m *LimitMiddleware) Handler(r rate.Limit, b int, t time.Duration) gin.HandlerFunc {
	limiters := &sync.Map{}

	return func(c *gin.Context) {
		// 获取限速器
		// key 除了 ip 之外也可以是其他的，例如 header，user name 等
		key := c.Request.Header.Get("X-Forward-For")
		// key = c.ClientIP()
		l, _ := limiters.LoadOrStore(key, rate.NewLimiter(r, b))

		// 这里注意不要直接使用 gin 的 context 默认是没有超时时间的
		ctx, cancel := context.WithTimeout(c, t)
		defer cancel()

		if err := l.(*rate.Limiter).Wait(ctx); err != nil {
			// 这里先不处理日志了，如果返回错误就直接 429
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": err})
		}
		c.Next()
	}
}

func NewLimitMiddleware() *LimitMiddleware {
	return &LimitMiddleware{}
}
