package gf

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/ulule/limiter/v3"
)

// Middleware is the middleware for goframe.
type Middleware struct {
	Limiter        *limiter.Limiter
	OnError        ErrorHandler
	OnLimitReached LimitReachedHandler
	KeyGetter      KeyGetter
	ExcludedKey    func(string) bool
}

// NewMiddleware return a new instance of a goframe middleware.
func NewMiddleware(limiter *limiter.Limiter, limitReachedHandler func(r *ghttp.Request)) ghttp.HandlerFunc {
	middleware := &Middleware{
		Limiter:        limiter,
		OnError:        DefaultErrorHandler,
		OnLimitReached: DefaultLimitReachedHandler,
		KeyGetter:      DefaultKeyGetter,
		ExcludedKey:    nil,
	}
	if limitReachedHandler != nil {
		middleware.OnLimitReached = limitReachedHandler
	}

	return func(r *ghttp.Request) {
		middleware.Handle(r)
	}
}

// Handle request.
func (middleware *Middleware) Handle(r *ghttp.Request) {
	key := middleware.KeyGetter(r)
	if middleware.ExcludedKey != nil && middleware.ExcludedKey(key) {
		r.Middleware.Next()
		return
	}

	ctx, err := middleware.Limiter.Get(r.Context(), key)
	if err != nil {
		middleware.OnError(r, err)
		r.ExitAll()
		return
	}

	if ctx.Reached {
		middleware.OnLimitReached(r)
		r.ExitAll()
		return
	}

	r.Middleware.Next()
}
