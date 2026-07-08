package logger

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

// Level defines the logging severity.
type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// LoggerPlugin handles routing Go logs to the native platform logger.
type LoggerPlugin struct {
	app *core.Application
	tag string
}

// NewPlugin creates a new instance of the LoggerPlugin.
func NewPlugin(tag string) *LoggerPlugin {
	if tag == "" {
		tag = "WailsGo"
	}
	return &LoggerPlugin{tag: tag}
}

// Name returns the plugin name.
func (p *LoggerPlugin) Name() string {
	return "logger"
}

// Init initializes the plugin.
func (p *LoggerPlugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

func (p *LoggerPlugin) log(level Level, message string) {
	payload, _ := json.Marshal(map[string]string{
		"tag":     p.tag,
		"level":   string(level),
		"message": message,
	})
	core.CallNativePlatform("logger:log", string(payload))
}

// Debug logs a debug message.
func (p *LoggerPlugin) Debug(format string, a ...interface{}) {
	p.log(LevelDebug, fmt.Sprintf(format, a...))
}

// Info logs an info message.
func (p *LoggerPlugin) Info(format string, a ...interface{}) {
	p.log(LevelInfo, fmt.Sprintf(format, a...))
}

// Warn logs a warning message.
func (p *LoggerPlugin) Warn(format string, a ...interface{}) {
	p.log(LevelWarn, fmt.Sprintf(format, a...))
}

// Error logs an error message.
func (p *LoggerPlugin) Error(format string, a ...interface{}) {
	p.log(LevelError, fmt.Sprintf(format, a...))
}
