package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func SaveEvaluationResult(rule bool, normal, anomaly float64, filePath string) error {
	// Open the file in append mode or create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Optionally, write headers if the file is new (check file size)
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		writer.Write([]string{"Rule", "Normal", "Anomaly"})
	}

	// Convert float values to strings
	record := []string{
		fmt.Sprint(rule),
		strconv.FormatFloat(normal, 'f', 6, 64),
		strconv.FormatFloat(anomaly, 'f', 6, 64),
	}

	// Write the record
	return writer.Write(record)
}
