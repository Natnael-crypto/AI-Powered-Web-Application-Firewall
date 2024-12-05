package logger

import (
	"fmt"
	"log"
	"os"
	"time"
	"ai-powered-waf/rule-engine/internal/config"
)

const (
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

var (
	logFile    *os.File
	debugMode  bool
	logChannel = make(chan string, 100)
)

// InitializeLogger sets up the logger, including file output and log level control.
func InitializeLogger() {
	var err error

	// Create or append to the log file
	logFile, err = os.OpenFile("./logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	debugMode = config.Config.Server.DebugMode

	// Start a goroutine to handle asynchronous logging
	go func() {
		for logMsg := range logChannel {
			fmt.Println(logMsg) // Print to console
			logFile.WriteString(logMsg + "\n")
		}
	}()
}

// writeLogEntry writes a formatted writeLogEntry message to the channel.
func writeLogEntry(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, level, message)

	logChannel <- logMsg
}

// Info logs an informational message.
func Info(message string) {
	writeLogEntry(INFO, message)
}

// Warn logs a warning message.
func Warn(message string) {
	writeLogEntry(WARN, message)
}

// Error logs an error message.
func Error(message string) {
	writeLogEntry(ERROR, message)
}

// LogRequest logs request details in debug mode.
func LogRequest(details string) {
	if debugMode {
		writeLogEntry(INFO, fmt.Sprintf("Request Details: %s", details))
	}
}

// CloseLogger closes the log file and channel.
func CloseLogger() {
	close(logChannel)
	logFile.Close()
}
