package main

import "fmt"

// 13. Adapter Pattern - Adapts incompatible interfaces to work together

// Target interface our client expects
type LogWriter interface {
	WriteLog(message string)
}

// External library with different interface
type ExternalLogger struct{}

func (e *ExternalLogger) Log(msg string) {
	fmt.Printf("[EXTERNAL] %s\n", msg)
}

// Adapter wraps ExternalLogger to implement LogWriter
type LoggerAdapter struct {
	external *ExternalLogger
}

func (a *LoggerAdapter) WriteLog(message string) {
	a.external.Log(message)
}

func NewLoggerAdapter(logger *ExternalLogger) *LoggerAdapter {
	return &LoggerAdapter{external: logger}
}

// Application that only knows LogWriter
func processWithLogger(lw LogWriter, msg string) {
	lw.WriteLog(msg)
}

func main() {
	extLogger := &ExternalLogger{}
	adapter := NewLoggerAdapter(extLogger)

	processWithLogger(adapter, "Application started")
	processWithLogger(adapter, "Request processed")
	processWithLogger(adapter, "Application stopped")
}
