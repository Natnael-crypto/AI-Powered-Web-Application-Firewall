package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RackSec/srslog"
)

var (
	logFile     *os.File
	syslogWriter *srslog.Writer
)

// InitializeLogger initializes both file and Syslog logging
func InitializeLogger(syslogAddr string) error {
	var err error

	// Open local file for logging
	logFile, err = os.OpenFile("waf.log", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to initialize file logger: %v", err)
	}

	// Initialize Syslog writer
	syslogWriter, err = srslog.Dial("udp", syslogAddr, srslog.LOG_INFO, "WAF")
	fmt.Print(syslogAddr)
	if err != nil {
		return fmt.Errorf("failed to initialize Syslog writer: %v", err)
	}

	log.Println("Logger initialized successfully.")
	return nil
}

// LogRequest logs HTTP requests to both Syslog and a local file
func LogRequest(r *http.Request, decision, reason string) {
	// Format log message
	message := fmt.Sprintf("Method: %s, URL: %s, Decision: %s, Reason: %s",
		r.Method, r.URL.String(), decision, reason)

	// Log to Syslog
	if syslogWriter != nil {
		err := syslogWriter.Info(message)
		if err != nil {
			log.Println("An error occured while trying to write logs:", err)
		}
	} else {
		log.Println("Syslog writer not initialized. Skipping Syslog logging.")
	}

	// Log to local file
	if logFile != nil {
		log.SetOutput(logFile)
		log.Println(message)
	} else {
		log.Println("File logger not initialized. Skipping file logging.")
	}
}

// CloseLogger closes the file and Syslog connections
func CloseLogger() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			log.Println("An error occurred trying to close log file:", err)
		}
	}
	if syslogWriter != nil {
		err := syslogWriter.Close()
		if err != nil {
			log.Println("An error occurred trying to close syslog writer:", err)
		}
	}
}
