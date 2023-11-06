package utils

import (
	"log/slog"
	"os"
	"time"

	"gitlab.com/greyxor/slogor"
)

func SetLogger(debug bool) {
	logLevel := slog.LevelInfo
	showSource := false
	if debug {
		logLevel = slog.LevelDebug
	}

	if os.Getenv("SC_DEBUG_SRC") == "true" {
		showSource = true
	}

	slog.SetDefault(slog.New(slogor.NewHandler(os.Stderr, &slogor.Options{
		TimeFormat: time.Stamp,
		Level:      logLevel,
		ShowSource: showSource,
	})))
}
