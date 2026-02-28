package main

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string, err error)
}
type SlogLogger struct {
	l *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	return &SlogLogger{
		l: logger,
	}
}

func (sl *SlogLogger) Info(msg string) {
	sl.l.Info(msg)
}

func (sl *SlogLogger) Error(msg string) {
	sl.l.Error(msg)
}

func (sl *SlogLogger) Fatal(msg string, err error) {
	sl.l.Error(fmt.Sprintf("%s: %v", msg, err))
	os.Exit(1)
}
