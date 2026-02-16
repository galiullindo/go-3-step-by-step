package main

import "log/slog"

func LogUserAction(logger *slog.Logger, user string, action string) {
	logger.Info("user action", slog.String("user", user), slog.String("action", action))
}
