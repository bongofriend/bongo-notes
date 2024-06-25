package api

import (
	"log"
	"net/http"
	"time"
)

type middleware func(http.Handler) http.Handler

type responseWithStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWithStatusCode) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.statusCode = status
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rsp := &responseWithStatusCode{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		h.ServeHTTP(rsp, r)
		log.Println(rsp.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func CreateMiddlewareStack(middlewares ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			m := middlewares[i]
			next = m(next)
		}
		return next
	}
}
