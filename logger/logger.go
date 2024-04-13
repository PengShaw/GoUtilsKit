package logger

import (
	"log"
	"os"
)

//go:generate stringer -type=LogLevel -trimprefix=Level

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

// A Logger extend log.Logger with LogLevel.
type Logger struct {
	*log.Logger
	level LogLevel
}

// New creates a new *[Logger].
func New(level LogLevel) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

var std = New(LevelInfo)

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }

// SetLevel sets the output level for the logger.
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Level returns the output level for the logger.
func (l *Logger) Level() string {
	return l.level.String()
}

// SetLevel sets the output level for the standard logger.
func SetLevel(level LogLevel) {
	std.level = level
}

// Level returns the output level for the standard logger.
func Level() string {
	return std.level.String()
}

// print by LogLevel

func (l *Logger) printf(level LogLevel, format string, v ...any) {
	if l.level <= level {
		l.Logger.Printf(format, v...)
	}
}

func (l *Logger) println(level LogLevel, v ...any) {
	if l.level <= level {
		l.Logger.Println(v...)
	}
}

func (l *Logger) print(level LogLevel, v ...any) {
	if l.level <= level {
		l.Logger.Print(v...)
	}
}

//go:generate go run gen.go
