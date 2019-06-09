package middleware

import (
	"data_base/presentation/logger"
	"net/http"
	"time"
)

func ContentType(this http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		this.ServeHTTP(w, r)
	})
}

func Logger(this http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		this.ServeHTTP(w, r)
		logger.Info.Printf("\nmethod: [%s]\naddr: %s\nparams: %s\npath: %s\nwork time: %s\n",
			r.Method, r.RemoteAddr, r.URL.RawQuery, r.URL.Path, time.Since(start))
	})
}

func Panic(this http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Fatal.Println("recovered", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		this.ServeHTTP(w, r)
	})
}
