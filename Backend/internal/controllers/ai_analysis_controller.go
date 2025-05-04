package controllers

import (
	"net/http"
	"sync"
	"time"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

var (
	analysisQueue = make(map[string]bool)
	queueMutex    sync.Mutex
)

func QueueRequestForAnalysis(c *gin.Context) {
	var input struct {
		RequestIDs []string `json:"request_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || len(input.RequestIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request_ids are required"})
		return
	}

	queueMutex.Lock()
	for _, reqID := range input.RequestIDs {
		analysisQueue[reqID] = true
	}
	queueMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Requests queued for AI analysis"})
}

// ML server fetches all queued requests and clears the queue
func FetchAndAnalyzeRequests(c *gin.Context) {
	queueMutex.Lock()
	var requestIDs []string
	for id := range analysisQueue {
		requestIDs = append(requestIDs, id)
	}
	// Clear the queue after fetching
	analysisQueue = make(map[string]bool)
	queueMutex.Unlock()

	if len(requestIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No requests in queue"})
		return
	}

	var requests []models.Request
	if err := config.DB.Where("request_id IN ?", requestIDs).Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

// ML server submits analysis results
func SubmitAnalysisResults(c *gin.Context) {
	var results []struct {
		RequestID  string `json:"request_id"`
		ThreatType string `json:"threat_type"`
	}

	if err := c.ShouldBindJSON(&results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	for _, result := range results {
		config.DB.Model(&models.Request{}).
			Where("request_id = ?", result.RequestID).
			Updates(map[string]interface{}{
				"ai_result":      true,
				"ai_threat_type": result.ThreatType,
			})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Analysis results updated"})
}

func CreateModelTrainingRequest(c *gin.Context) {
	var input struct {
		ID                 string  `json:"id" binding:"required"`
		NumberRequestsUsed int     `json:"number_requests_used" binding:"required"`
		PercentTrainData   float32 `json:"percent_train_data" binding:"required"`
		NumTrees           int     `json:"num_trees" binding:"required"`
		MaxDepth           int     `json:"max_depth" binding:"required"`
		MinSamplesSplit    int     `json:"min_samples_split" binding:"required"`
		Criterion          string  `json:"criterion" binding:"required"` // "gini" or "entropy"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	ai_model := models.AIModel{
		ID:                 input.ID,
		NumberRequestsUsed: input.NumberRequestsUsed,
		PercentTrainData:   input.PercentTrainData,
		NumTrees:           input.NumTrees,
		MaxDepth:           input.MaxDepth,
		MinSamplesSplit:    input.MinSamplesSplit,
		Criterion:          input.Criterion,
		Accuracy:           0,
		Precision:          0,
		Recall:             0,
		F1:                 0,
		Selected:           false,
		Modeled:            false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := config.DB.Create(&ai_model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create model"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "model training request created", "model": ai_model})
}

func GetUntrainedModelForML(c *gin.Context) {
	var model models.AIModel

	if err := config.DB.Where("modeled = ?", false).First(&model).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "no untrained model available"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"model": model})
}

func SubmitModelResults(c *gin.Context) {
	var input struct {
		ID        string  `json:"id"`
		Accuracy  float32 `json:"accuracy"`
		Precision float32 `json:"precision"`
		F1        float32 `json:"f1"`
		Recall    float32 `json:"recall"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := config.DB.Model(&models.AIModel{}).Where("id = ?", input.ID).Updates(map[string]interface{}{
		"accuracy":  input.Accuracy,
		"precision": input.Precision,
		"recall":    input.Recall,
		"fl":        input.F1,
		"modeled":   true,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update model stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "model results updated"})
}

func SelectActiveModel(c *gin.Context) {
	var input struct {
		ID string `json:"id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := config.DB.Model(&models.AIModel{}).Where("selected = ?", true).Update("selected", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deselect existing model"})
		return
	}

	if err := config.DB.Model(&models.AIModel{}).Where("id = ?", input.ID).Update("selected", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to select model"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "model selected successfully"})
}

func GetSelectedModel(c *gin.Context) {
	var model models.AIModel

	if err := config.DB.Where("selected = ?", true).First(&model).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "no model currently selected"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"model": model})
}

func GetModels(c *gin.Context) {
	var model models.AIModel

	if err := config.DB.Find(&model).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "no model currently selected"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"model": model})
}
