package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(next http.Handler) http.Handler {
	logger := slog.Default()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(lw, r)

		logStr := fmt.Sprintf("%s,%s,%d,%s\n", r.Method, r.URL.Path, lw.status, time.Since(start))

		logger.Info("Request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lw.status,
			"latency", time.Since(start),
		)

		// TODO: use channels
		go func(logMessage string) {
			_, err := writeLog(logMessage)
			if err != nil {
				slog.Error("Failed to write log asynchronously", "error", err)
			}
		}(logStr)
	})
}

func writeLog(msg string) (int, error) {
	f, err := os.OpenFile("logs.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		return 0, err
	}

	defer f.Close()

	n, err := f.Write([]byte(msg))
	if err != nil {
		return 0, fmt.Errorf("error occurred while writing: %w", err)
	}

	return n, nil
}
