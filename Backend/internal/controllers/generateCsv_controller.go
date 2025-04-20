package controllers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateRequestsCSV handles the generation and download of requests data as a CSV file
func GenerateRequestsCSV(c *gin.Context) {

	var requests []models.Request
	query := utils.ApplyRequestFilters(c)
	if query == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to apply filters"})
		return
	}

	if err := query.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch requests"})
		return
	}

	// Set headers for CSV download
	filename := fmt.Sprintf("requests_%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")

	// Create a CSV writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write CSV header
	headers := []string{
		"ID", "Application Name", "Client IP", "Request Method", "Request URL", "Headers", "Body", "Timestamp", "ResponseCode", "Status", "ThreatType", "BotDetected", "GeoLocation", "RateLimited", "UserAgent",
		// Add more fields as needed based on your request structure
	}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write CSV headers",
		})
		return
	}

	// Write request data rows
	for _, req := range requests {
		row := []string{
			req.RequestID,
			req.ApplicationName,
			req.ClientIP,
			req.RequestMethod,
			req.RequestURL,
			req.Headers,
			req.Body,
			fmt.Sprintf("%.3f", req.Timestamp),
			fmt.Sprintf("%d", req.ResponseCode),
			req.Status,
			req.ThreatType,
			fmt.Sprintf("%t", req.BotDetected),
			req.GeoLocation,
			fmt.Sprintf("%t", req.RateLimited),
			req.UserAgent,
			// Add more fields as needed
		}
		if err := writer.Write(row); err != nil {
			// Log the error but continue processing
			fmt.Printf("Error writing request %s to CSV: %v\n", req.RequestID, err)
		}
	}
	// The writer will write to the response when Flush() is called
}
