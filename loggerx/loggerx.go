package loggerx

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type loggerOpts struct {
	consolePretty bool
	level         slog.Level
	w             io.Writer
}

type LoggerOption interface {
	apply(*loggerOpts)
}

type LoggerOptions []LoggerOption

func (o LoggerOptions) apply(s *loggerOpts) {
	for _, opt := range o {
		opt.apply(s)
	}
}

type consolePrettyOpt bool

func (o consolePrettyOpt) apply(s *loggerOpts) { s.consolePretty = bool(o) }
func WithConsolePretty(flag bool) LoggerOption { return consolePrettyOpt(flag) }

type levelOpt slog.Level

func (o levelOpt) apply(s *loggerOpts)        { s.level = slog.Level(o) }
func WithLevel(level slog.Level) LoggerOption { return levelOpt(level) }

type writerOpt struct{ w io.Writer }

func (o writerOpt) apply(s *loggerOpts)   { s.w = o.w }
func WithWriter(w io.Writer) LoggerOption { return writerOpt{w: w} }

func New(opts ...LoggerOption) *slog.Logger {
	lopts := &loggerOpts{
		w:             os.Stdout,
		consolePretty: true,
		level:         slog.LevelDebug,
	}

	for _, o := range opts {
		o.apply(lopts)
	}

	jsonHandler := slog.NewJSONHandler(lopts.w, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})

	prettyHandler := tint.NewHandler(lopts.w, &tint.Options{
		Level:      lopts.level,
		AddSource:  false,
		TimeFormat: time.TimeOnly,
	})

	var handler slog.Handler = jsonHandler
	if lopts.consolePretty {
		handler = prettyHandler
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
