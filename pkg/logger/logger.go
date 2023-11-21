package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger extends logrus.Logger with additional functionality
type Logger struct {
	*logrus.Logger
	fields   logrus.Fields
	name     string
	parent   *Logger
	children map[string]*Logger
}

// Config holds logger configuration
type Config struct {
	Level        string
	ReportCaller bool
	JSONFormat   bool
	FileOutput   string
	TimeFormat   string
	TreeFormat   bool
	UseColors    bool
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:        "info",
		ReportCaller: true,
		JSONFormat:   false,
		TimeFormat:   time.RFC3339,
	}
}

// New creates a new logger instance with given configuration
func New(config *Config) (*Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Create base logrus logger
	log := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	log.SetLevel(level)

	// Configure formatter
	if config.TreeFormat {
		log.SetFormatter(&TreeFormatter{
			TimestampFormat: config.TimeFormat,
			ShowCaller:      config.ReportCaller,
			UseColors:       config.UseColors,
		})
	} else if config.JSONFormat {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: config.TimeFormat,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := filepath.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", s, f.Line)
			},
		})
	} else {
		formatter := &logrus.TextFormatter{
			TimestampFormat: config.TimeFormat,
			FullTimestamp:   true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
		log.SetFormatter(formatter)
	}

	// Configure output
	if config.FileOutput != "" {
		file, err := os.OpenFile(config.FileOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		log.SetOutput(file)
	}

	// Enable caller reporting if configured
	log.SetReportCaller(config.ReportCaller)

	return &Logger{
		Logger: log,
		fields: logrus.Fields{},
	}, nil
}

// WithField adds a field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	return &Logger{
		Logger: l.Logger,
		fields: newFields,
	}
}

// WithFields adds multiple fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	return &Logger{
		Logger: l.Logger,
		fields: newFields,
	}
}

// WithError adds an error to the logger context
func (l *Logger) WithError(err error) *Logger {
	return l.WithField("error", err)
}

// log implements the actual logging logic
func (l *Logger) log(level logrus.Level, args ...interface{}) {
	if len(l.fields) > 0 {
		l.Logger.WithFields(l.fields).Log(level, args...)
	} else {
		l.Logger.Log(level, args...)
	}
}

// logf implements the actual formatted logging logic
func (l *Logger) logf(level logrus.Level, format string, args ...interface{}) {
	if len(l.fields) > 0 {
		l.Logger.WithFields(l.fields).Logf(level, format, args...)
	} else {
		l.Logger.Logf(level, format, args...)
	}
}

// Convenience methods for different log levels
func (l *Logger) Debug(args ...interface{}) { l.log(logrus.DebugLevel, args...) }
func (l *Logger) Info(args ...interface{})  { l.log(logrus.InfoLevel, args...) }
func (l *Logger) Warn(args ...interface{})  { l.log(logrus.WarnLevel, args...) }
func (l *Logger) Error(args ...interface{}) { l.log(logrus.ErrorLevel, args...) }
func (l *Logger) Fatal(args ...interface{}) { l.log(logrus.FatalLevel, args...) }
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(logrus.DebugLevel, format, args...)
}
func (l *Logger) Infof(format string, args ...interface{}) { l.logf(logrus.InfoLevel, format, args...) }
func (l *Logger) Warnf(format string, args ...interface{}) { l.logf(logrus.WarnLevel, format, args...) }
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(logrus.ErrorLevel, format, args...)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(logrus.FatalLevel, format, args...)
}

// SubLoggerOpts contains options for creating a sub-logger
type SubLoggerOpts struct {
	// Additional fields to add to the sub-logger
	Fields logrus.Fields
	// Override the log level for this sub-logger (optional)
	Level *logrus.Level
}

// NewSubLogger creates a new sub-logger with the given name
func (l *Logger) NewSubLogger(name string, opts *SubLoggerOpts) *Logger {
	if opts == nil {
		opts = &SubLoggerOpts{}
	}

	// Merge parent fields with new fields
	fields := make(logrus.Fields)
	for k, v := range l.fields {
		fields[k] = v
	}
	if opts.Fields != nil {
		for k, v := range opts.Fields {
			fields[k] = v
		}
	}

	// Always add the logger name to fields
	fields["logger"] = name
	if l.name != "" {
		fields["logger"] = l.name + "." + name
	}

	subLogger := &Logger{
		Logger:   l.Logger,
		fields:   fields,
		name:     name,
		parent:   l,
		children: make(map[string]*Logger),
	}

	// Store in parent's children map
	if l.children == nil {
		l.children = make(map[string]*Logger)
	}
	l.children[name] = subLogger

	return subLogger
}

// GetSubLogger retrieves an existing sub-logger by name
func (l *Logger) GetSubLogger(name string) *Logger {
	if l.children == nil {
		return nil
	}
	return l.children[name]
}

// GetAllSubLoggers returns all immediate sub-loggers
func (l *Logger) GetAllSubLoggers() map[string]*Logger {
	return l.children
}

// WithScope adds a scope field to the logger
func (l *Logger) WithScope(scope string) *Logger {
	return l.WithField("scope", scope)
}

// WithComponent adds a component field to the logger
func (l *Logger) WithComponent(component string) *Logger {
	return l.WithField("component", component)
}
