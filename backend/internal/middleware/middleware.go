package middleware

import (
	"context"
	"mini-app-backend/internal/logger"
	"mini-app-backend/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)
		
		w.Header().Set("X-Request-ID", requestID)
		
		logger.WithRequestID(requestID).Infof("Started %s %s", r.Method, r.URL.Path)
		
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
		
		status := rw.status
		
		logger.WithRequestID(requestID).Infof("Completed %s %s in %v. Status: %d", r.Method, r.URL.Path, time.Since(start), status)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		if origin == "http://localhost:3030" || origin == "https://7wt1l5rz-3030.euw.devtunnels.ms" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err == nil {
			encryptionUtil := utils.NewEncryptionUtil()
			
			decryptedValue, err := encryptionUtil.Decrypt(cookie.Value)
			if err == nil {
				userID, err := strconv.ParseInt(decryptedValue, 10, 64)
				if err == nil {
					ctx := context.WithValue(r.Context(), "user_id", userID)
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
					return
				}
			}
			
			userID, err := strconv.ParseInt(cookie.Value, 10, 64)
			if err == nil {
				ctx := context.WithValue(r.Context(), "user_id", userID)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}