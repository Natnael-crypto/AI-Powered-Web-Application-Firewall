package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RackSec/srslog"
)

var (
	logFile      *os.File
	syslogWriter *srslog.Writer
)

func InitializeLogger(syslogAddr string) error {
	var err error

	logFile, err = os.OpenFile("waf.log", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to initialize file logger: %v", err)
	}

	syslogWriter, err = srslog.Dial("udp", syslogAddr, srslog.LOG_INFO, "WAF")
	fmt.Print(syslogAddr)
	if err != nil {
		return fmt.Errorf("failed to initialize Syslog writer: %v", err)
	}

	log.Println("Logger initialized successfully.")
	return nil
}

func LogRequest(r *http.Request, decision string, application string, ip string, anomalyScore float64) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	severity := 6

	if anomalyScore > 0.5 {
		severity = 2
	}

	PRI := 20*8 + severity
	syslogHeader := fmt.Sprintf("<%d>1 %s %s %s - %s -", PRI,
		timestamp, hostname, "Gasha waf", application)

	jsonBody := fmt.Sprintf(`{"method":"%s","url":"%s","decision":"%s","source_ip":"%s","anomaly_score":%.2f}`,
		r.Method, r.URL.String(), decision, ip, anomalyScore)

	fullMessage := fmt.Sprintf("%s %s", syslogHeader, jsonBody)

	if syslogWriter != nil {
		err := syslogWriter.Info(fullMessage)
		if err != nil {
			log.Println("An error occurred while trying to write to syslog:", err)
		}
	} else {
		log.Println("Syslog writer not initialized. Skipping syslog logging.")
	}

	if logFile != nil {
		log.SetOutput(logFile)
		log.Println(fullMessage)
	} else {
		log.Println("File logger not initialized. Skipping file logging.")
	}
}

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
