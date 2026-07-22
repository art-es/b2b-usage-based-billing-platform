package log

import (
	"io"
	"os"
	"time"
)

const (
	Disabled Level = iota + 1
	Error
	Warning
	Info
	Debug
)

const (
	FormatJSON Format = iota + 1
	FormatText
)

var (
	defaultLevel        = Info
	defaultFormat       = FormatJSON
	defaultOutput       = os.Stdout
	defaultGetCreatedAt = func() string { return time.Now().Format(time.DateTime) }
)

type Level int8

type Format int8

type Logger interface {
	Set(key, val string) Logger
	Log(lvl Level) Log
}

type Log interface {
	Set(key, val string) Log
	Write()
}

type Options struct {
	// Log level. Default: Info
	Level Level

	// Log format. Default: JSON
	Format Format

	// The target where to send the generated log. Default: os.Stdout
	Output io.Writer

	// Returns a value for "created_at" field. Default: time.Now().Format(time.DateTime)
	GetCreatedAt func() string
}

func NewLogger(opts *Options) Logger {
	if opts == nil {
		opts = &Options{}
	}

	if !opts.Level.IsValid() {
		opts.Level = defaultLevel
	}

	if !opts.Format.IsValid() {
		opts.Format = defaultFormat
	}

	if opts.Output == nil {
		opts.Output = defaultOutput
	}

	if opts.GetCreatedAt == nil {
		opts.GetCreatedAt = defaultGetCreatedAt
	}

	return logger{
		level:        opts.Level,
		format:       opts.Format,
		output:       opts.Output,
		getCreatedAt: opts.GetCreatedAt,
	}
}

func (l Level) String() string {
	switch l {
	case Disabled:
		return "disabled"
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Info:
		return "info"
	case Debug:
		return "debug"
	default:
		return "unknown"
	}
}

func (l Level) IsValid() bool {
	switch l {
	case Disabled, Error, Warning, Info, Debug:
		return true
	default:
		return false
	}
}

func (f Format) IsValid() bool {
	switch f {
	case FormatJSON, FormatText:
		return true
	default:
		return false
	}
}
