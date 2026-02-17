package main

import "log/slog"

func LogHTTPRequest(logger *slog.Logger, method string, path string, status int, duration_ms int64) {
	logger.Info(
		"http request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", status),
		slog.Int64("duration_ms", duration_ms),
	)
}
