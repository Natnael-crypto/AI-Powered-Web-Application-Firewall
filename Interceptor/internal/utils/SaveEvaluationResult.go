package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func SaveEvaluationResult(rule bool, normal, anomaly float64, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		writer.Write([]string{"Rule", "Normal", "Anomaly"})
	}

	record := []string{
		fmt.Sprint(rule),
		strconv.FormatFloat(normal, 'f', 6, 64),
		strconv.FormatFloat(anomaly, 'f', 6, 64),
	}

	return writer.Write(record)
}
