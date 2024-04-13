package logger_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/PengShaw/GoUtilsKit/logger"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestLogLevelString(t *testing.T) {
	gots := []logger.LogLevel{
		logger.LevelTrace,
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
		logger.LevelPanic,
		logger.LogLevel(7),
	}
	wants := []string{
		"Trace",
		"Debug",
		"Info",
		"Warn",
		"Error",
		"Fatal",
		"Panic",
		"LogLevel(7)",
	}
	for i := range wants {
		got := gots[i]
		want := wants[i]
		assert.Equal(t, want, got.String(), "they should be equal")
	}
}

func TestLogPanic(t *testing.T) {
	t.Run("test a new logger", func(t *testing.T) {
		l := logger.New(logger.LevelPanic)
		var got bytes.Buffer
		l.SetOutput(&got)
		assert.Equal(t, "Panic", l.Level(), "they should be equal")
		// Panicf, Panicln, Panic
		assert.PanicsWithValue(t, "[PANIC] test Panicf with msg", func() { l.Panicf("test Panicf %s", "with msg") })
		assert.PanicsWithValue(t, "[PANIC] test Panicln\n", func() { l.Panicln("test Panicln") })
		assert.PanicsWithValue(t, "[PANIC] test Panic", func() { l.Panic("test Panic") })
		assert.Contains(t, got.String(), "[PANIC] test Panicf with msg\n")
		assert.Contains(t, got.String(), "[PANIC] test Panicln\n")
		assert.Contains(t, got.String(), "[PANIC] test Panic\n")
		// Fatal should not be execute, because the LogLevel
		l.Fatalln("Fatalln")
		assert.NotContains(t, got.String(), "[FATAL] Fatalln\n")
	})

	t.Run("test std logger", func(t *testing.T) {
		l := logger.Default()
		l.SetLevel(logger.LevelPanic)
		var got bytes.Buffer
		l.SetOutput(&got)
		assert.Equal(t, "Panic", logger.Level(), "they should be equal")
		// Panicf, Panicln, Panic
		assert.PanicsWithValue(t, "[PANIC] test Panicf with msg", func() { logger.Panicf("test Panicf %s", "with msg") })
		assert.PanicsWithValue(t, "[PANIC] test Panicln\n", func() { logger.Panicln("test Panicln") })
		assert.PanicsWithValue(t, "[PANIC] test Panic", func() { logger.Panic("test Panic") })
		assert.Contains(t, got.String(), "[PANIC] test Panicf with msg\n")
		assert.Contains(t, got.String(), "[PANIC] test Panicln\n")
		assert.Contains(t, got.String(), "[PANIC] test Panic\n")
		// Fatal should not be execute, because the LogLevel
		logger.Fatalln("Fatalln")
		assert.NotContains(t, got.String(), "Fatalln\n")
	})
}

func TestLogFatal(t *testing.T) {

	t.Run("test a new logger", func(t *testing.T) {
		count := 0
		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			count += 1
		})
		defer patches.Reset()

		l := logger.New(logger.LevelFatal)
		var got bytes.Buffer
		l.SetOutput(&got)
		assert.Equal(t, "Fatal", l.Level(), "they should be equal")
		// Fatalf, Fatalln, Fatal
		l.Fatalf("test Fatalf %s", "with msg")
		l.Fatalln("test Fatalln")
		l.Fatal("test Fatal")
		assert.Contains(t, got.String(), "[FATAL] test Fatalf with msg\n")
		assert.Contains(t, got.String(), "[FATAL] test Fatalln\n")
		assert.Contains(t, got.String(), "[FATAL] test Fatal\n")
		assert.Equal(t, 3, count, "Fatal times should be except")
		// Error should not be execute, because the LogLevel
		l.Errorln("Errorln")
		assert.NotContains(t, got.String(), "Errorln\n")
	})

	t.Run("test std logger", func(t *testing.T) {
		count := 0
		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			count += 1
		})
		defer patches.Reset()

		l := logger.Default()
		var got bytes.Buffer
		l.SetOutput(&got)
		l.SetLevel(logger.LevelFatal)

		assert.Equal(t, "Fatal", logger.Level(), "they should be equal")
		// Fatalf, Fatalln, Fatal
		logger.Fatalf("test Fatalf %s", "with msg")
		logger.Fatalln("test Fatalln")
		logger.Fatal("test Fatal")
		assert.Contains(t, got.String(), "[FATAL] test Fatalf with msg\n")
		assert.Contains(t, got.String(), "[FATAL] test Fatalln\n")
		assert.Contains(t, got.String(), "[FATAL] test Fatal\n")
		assert.Equal(t, 3, count, "Fatal times should be except")
		// Error should not be execute, because the LogLevel
		logger.Errorln("Errorln")
		assert.NotContains(t, got.String(), "Errorln\n")
	})
}

func TestLogError(t *testing.T) {
	t.Run("test a new logger", func(t *testing.T) {
		l := logger.New(logger.LevelError)
		var got bytes.Buffer
		l.SetOutput(&got)
		assert.Equal(t, "Error", l.Level(), "they should be equal")
		// Errorf, Errorln, Error
		l.Errorf("test Errorf %s", "with msg")
		l.Errorln("test Errorln")
		l.Error("test Error")
		assert.Contains(t, got.String(), "[ERROR] test Errorf with msg\n")
		assert.Contains(t, got.String(), "[ERROR] test Errorln\n")
		assert.Contains(t, got.String(), "[ERROR] test Error\n")
		// WARN should not be execute, because the LogLevel
		l.Warnln("Warnln")
		assert.NotContains(t, got.String(), "Warnln\n")
	})

	t.Run("test std logger", func(t *testing.T) {
		l := logger.Default()
		var got bytes.Buffer
		l.SetOutput(&got)
		l.SetLevel(logger.LevelError)
		assert.Equal(t, "Error", logger.Level(), "they should be equal")
		// Errorf, Errorln, Error
		logger.Errorf("test Errorf %s", "with msg")
		logger.Errorln("test Errorln")
		logger.Error("test Error")
		assert.Contains(t, got.String(), "[ERROR] test Errorf with msg\n")
		assert.Contains(t, got.String(), "[ERROR] test Errorln\n")
		assert.Contains(t, got.String(), "[ERROR] test Error\n")
		// WARN should not be execute, because the LogLevel
		logger.Warnln("Warnln")
		assert.NotContains(t, got.String(), "Warnln\n")
	})
}

func TestLogWarn(t *testing.T) {
	t.Run("test a new logger", func(t *testing.T) {
		l := logger.New(logger.LevelWarn)
		var got bytes.Buffer
		l.SetOutput(&got)
		assert.Equal(t, "Warn", l.Level(), "they should be equal")
		// Warnf, Warnln, Warn
		l.Warnf("test Warnf %s", "with msg")
		l.Warnln("test Warnln")
		l.Warn("test Warn")
		assert.Contains(t, got.String(), "[WARN] test Warnf with msg\n")
		assert.Contains(t, got.String(), "[WARN] test Warnln\n")
		assert.Contains(t, got.String(), "[WARN] test Warn\n")
		// Info should not be execute, because the LogLevel
		l.Infoln("Infoln")
		assert.NotContains(t, got.String(), "Infoln\n")
	})

	t.Run("test std logger", func(t *testing.T) {
		l := logger.Default()
		var got bytes.Buffer
		l.SetOutput(&got)
		l.SetLevel(logger.LevelWarn)
		assert.Equal(t, "Warn", logger.Level(), "they should be equal")
		// Warnf, Warnln, Warn
		logger.Warnf("test Warnf %s", "with msg")
		logger.Warnln("test Warnln")
		logger.Warn("test Warn")
		assert.Contains(t, got.String(), "[WARN] test Warnf with msg\n")
		assert.Contains(t, got.String(), "[WARN] test Warnln\n")
		assert.Contains(t, got.String(), "[WARN] test Warn\n")
		// Info should not be execute, because the LogLevel
		logger.Infoln("Infoln")
		assert.NotContains(t, got.String(), "Infoln\n")
	})
}

func TestLogInfo(t *testing.T) {
	t.Run("test a new logger", func(t *testing.T) {
		l := logger.New(logger.LevelInfo)
		var got bytes.Buffer
		l.SetOutput(&got)
		assert.Equal(t, "Info", l.Level(), "they should be equal")
		// Infof, Infoln, Info
		l.Infof("test Infof %s", "with msg")
		l.Infoln("test Infoln")
		l.Info("test Info")
		assert.Contains(t, got.String(), "[INFO] test Infof with msg\n")
		assert.Contains(t, got.String(), "[INFO] test Infoln\n")
		assert.Contains(t, got.String(), "[INFO] test Info\n")
		// Debug should not be execute, because the LogLevel
		l.Debugln("Debugln")
		assert.NotContains(t, got.String(), "Debugln\n")
	})

	t.Run("test std logger", func(t *testing.T) {
		l := logger.Default()
		var got bytes.Buffer
		l.SetOutput(&got)
		l.SetLevel(logger.LevelInfo)
		assert.Equal(t, "Info", logger.Level(), "they should be equal")
		// Infof, Infoln, Info
		logger.Infof("test Infof %s", "with msg")
		logger.Infoln("test Infoln")
		logger.Info("test Info")
		assert.Contains(t, got.String(), "[INFO] test Infof with msg\n")
		assert.Contains(t, got.String(), "[INFO] test Infoln\n")
		assert.Contains(t, got.String(), "[INFO] test Info\n")
		// Debug should not be execute, because the LogLevel
		logger.Debugln("Debugln")
		assert.NotContains(t, got.String(), "Debugln\n")
	})
}

func TestLogDebug(t *testing.T) {
	t.Run("test a new logger", func(t *testing.T) {
		l := logger.New(logger.LevelDebug)
		var got bytes.Buffer
		l.SetOutput(&got)
		assert.Equal(t, "Debug", l.Level(), "they should be equal")
		// Debugf, Debugln, Debug
		l.Debugf("test Debugf %s", "with msg")
		l.Debugln("test Debugln")
		l.Debug("test Debug")
		assert.Contains(t, got.String(), "[DEBUG] test Debugf with msg\n")
		assert.Contains(t, got.String(), "[DEBUG] test Debugln\n")
		assert.Contains(t, got.String(), "[DEBUG] test Debug\n")
		// Trace should not be execute, because the LogLevel
		l.Traceln("Traceln")
		assert.NotContains(t, got.String(), "Traceln\n")
	})

	t.Run("test std logger", func(t *testing.T) {
		l := logger.Default()
		var got bytes.Buffer
		l.SetOutput(&got)
		l.SetLevel(logger.LevelDebug)
		assert.Equal(t, "Debug", logger.Level(), "they should be equal")
		// Debugf, Debugln, Debug
		logger.Debugf("test Debugf %s", "with msg")
		logger.Debugln("test Debugln")
		logger.Debug("test Debug")
		assert.Contains(t, got.String(), "[DEBUG] test Debugf with msg\n")
		assert.Contains(t, got.String(), "[DEBUG] test Debugln\n")
		assert.Contains(t, got.String(), "[DEBUG] test Debug\n")
		// Trace should not be execute, because the LogLevel
		logger.Traceln("Traceln")
		assert.NotContains(t, got.String(), "Traceln\n")
	})
}
func TestLogTrace(t *testing.T) {
	t.Run("test a new logger", func(t *testing.T) {
		l := logger.New(logger.LevelTrace)
		var got bytes.Buffer
		l.SetOutput(&got)
		assert.Equal(t, "Trace", l.Level(), "they should be equal")
		// Tracef, Traceln, Trace
		l.Tracef("test Tracef %s", "with msg")
		l.Traceln("test Traceln")
		l.Trace("test Trace")
		assert.Contains(t, got.String(), "[TRACE] test Tracef with msg\n")
		assert.Contains(t, got.String(), "[TRACE] test Traceln\n")
		assert.Contains(t, got.String(), "[TRACE] test Trace\n")
	})

	t.Run("test std logger", func(t *testing.T) {
		l := logger.Default()
		var got bytes.Buffer
		l.SetOutput(&got)
		l.SetLevel(logger.LevelTrace)
		assert.Equal(t, "Trace", logger.Level(), "they should be equal")
		// Tracef, Traceln, Trace
		logger.Tracef("test Tracef %s", "with msg")
		logger.Traceln("test Traceln")
		logger.Trace("test Trace")
		assert.Contains(t, got.String(), "[TRACE] test Tracef with msg\n")
		assert.Contains(t, got.String(), "[TRACE] test Traceln\n")
		assert.Contains(t, got.String(), "[TRACE] test Trace\n")
	})
}
