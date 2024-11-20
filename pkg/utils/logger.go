package utils

import (
	"log/slog"
	"os"
)

func SetLogger(debug bool) {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}

	if debug {
		opts.Level = slog.LevelDebug
	}

	if os.Getenv("SC_DEBUG_SRC") == "true" {
		opts.AddSource = true
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
}
