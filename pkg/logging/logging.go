package logging

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var (
	log    *zap.Logger
	sugar  *zap.SugaredLogger
	once   sync.Once
	logMux sync.Mutex
)

// customEncoder creates a custom encoder that puts log level before timestamp
func customEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	return &customConsoleEncoder{zapcore.NewConsoleEncoder(config)}
}

// ValidLogLevels contains all valid log level options
var ValidLogLevels = []string{"DEBUG", "INFO", "WARN", "ERROR"}

// Initialize sets up the logger with the specified log level
func Initialize(logLevel string) error {
	var err error
	once.Do(func() {
		// Create a basic encoder configuration
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:          "time",
			LevelKey:         "level",
			NameKey:          "logger",
			CallerKey:        "caller",
			MessageKey:       "msg",
			StacktraceKey:    "stacktrace",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeDuration:   zapcore.SecondsDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: "\t",
		}

		// Use the custom encoder to change the order of fields (level first, then timestamp)

		// Use the custom encoder instead of the default console encoder
		consoleEncoder := customEncoder(encoderConfig)

		// Create a core that writes to stdout
		core := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			getZapLevel(logLevel),
		)

		// Create the logger
		log = zap.New(core)
		sugar = log.Sugar()
	})

	// If the logger is already initialized but we need to change the level
	if log != nil && err == nil {
		logMux.Lock()
		defer logMux.Unlock()

		// Validate log level
		level := getZapLevel(logLevel)
		if level == zapcore.InvalidLevel {
			return fmt.Errorf("invalid log level: %s. Valid options are: %s", logLevel, strings.Join(ValidLogLevels, ", "))
		}

		// Create a new logger with the updated level
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:          "time",
			LevelKey:         "level",
			NameKey:          "logger",
			CallerKey:        "caller",
			MessageKey:       "msg",
			StacktraceKey:    "stacktrace",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeDuration:   zapcore.SecondsDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: "\t",
		}

		// Use the custom encoder instead of the default console encoder
		consoleEncoder := customEncoder(encoderConfig)
		core := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)

		log = zap.New(core)
		sugar = log.Sugar()
	}

	return nil
}

// getZapLevel converts a string log level to a zap.AtomicLevel
func getZapLevel(logLevel string) zapcore.Level {
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InvalidLevel
	}
}

// Debug logs a debug message
func Debug(format string, args ...any) {
	if sugar != nil {
		sugar.Debugf(format, args...)
	}
}

// Info logs an info message
func Info(format string, args ...any) {
	if sugar != nil {
		sugar.Infof(format, args...)
	}
}

// Warn logs a warning message
func Warn(format string, args ...any) {
	if sugar != nil {
		sugar.Warnf(format, args...)
	}
}

// Error logs an error message
func Error(format string, args ...any) {
	if sugar != nil {
		sugar.Errorf(format, args...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(format string, args ...any) {
	if sugar != nil {
		sugar.Fatalf(format, args...)
	}
}

// Print logs a message that should always be shown regardless of verbosity
// This is implemented by temporarily setting the log level to INFO
func Print(format string, args ...any) {
	// Create a temporary logger with INFO level

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: "\t",
	}

	// Use the custom encoder
	consoleEncoder := customEncoder(encoderConfig)
	core := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	tempLogger := zap.New(core).Sugar()
	tempLogger.Infof(format, args...)
}

// customConsoleEncoder is a wrapper around the console encoder that changes the order of fields
type customConsoleEncoder struct {
	zapcore.Encoder
}

// Clone implements zapcore.Encoder
func (e *customConsoleEncoder) Clone() zapcore.Encoder {
	return &customConsoleEncoder{e.Encoder.Clone()}
}

// EncodeEntry implements zapcore.Encoder
func (e *customConsoleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// Get the original buffer
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// Get the buffer content as string
	content := buf.String()

	// Find the timestamp and level in the log line
	parts := strings.Split(content, "\t")
	if len(parts) >= 2 {
		// Swap the timestamp and level parts
		parts[0], parts[1] = parts[1], parts[0]

		// Create a new buffer with the modified content
		newBuf := buffer.NewPool().Get()
		newBuf.AppendString(strings.Join(parts, "\t"))
		return newBuf, nil
	}

	// If we couldn't parse the log line correctly, return the original
	return buf, nil
}
