package controllers

import (
	"net/http"
	"sync"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

var (
	analysisQueue = make(map[string]bool)
	queueMutex    sync.Mutex
)

func QueueRequestForAnalysis(c *gin.Context) {

	var input struct {
		RequestID string `json:"request_ids"`
	}

	var request models.Request

	if err := config.DB.Where("request_id = ?", input.RequestID).First(&request).Error; err != nil {
		c.JSON(http.StatusFound, gin.H{"error": "request not found"})
		return
	}

	var application models.Application

	if err := config.DB.Where("request_id = ?", request.ApplicationName).First(&application).Error; err != nil {
		c.JSON(http.StatusFound, gin.H{"error": "Application not found"})
		return
	}

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, application.ApplicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request_id are required"})
		return
	}

	queueMutex.Lock()
	analysisQueue[input.RequestID] = true
	queueMutex.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Request queued for AI analysis"})
}

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

func GetModelForMLs(c *gin.Context) {
	var models models.AIModel

	if err := config.DB.Find(&models).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "no models available"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"models": models})
}

func SubmitTrainResults(c *gin.Context) {
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
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update model stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "model results updated"})
}

func SelectActiveModel(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	modelID := c.Param("model_id")

	if err := config.DB.Where("modeled = ? AND id = ?", true, modelID).First(&models.AIModel{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "model not found or not modeled"})
		return
	}

	if err := config.DB.Model(&models.AIModel{}).Where("selected = ?", true).Update("selected", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deselect existing model"})
		return
	}

	if err := config.DB.Model(&models.AIModel{}).Where("id = ?", modelID).Update("selected", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to select model"})
		return
	}

	config.SelectModel = true

	c.JSON(http.StatusOK, gin.H{"message": "model selected successfully"})
}

func GetSelectedModel(c *gin.Context) {
	var model models.AIModel

	if err := config.DB.Where("selected = ?", true).First(&model).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "no model currently selected"})
		return
	}

	config.SelectModel = false

	c.JSON(http.StatusOK, gin.H{"model": model})
}

func GetModels(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var model []models.AIModel

	if err := config.DB.Find(&model).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "no model currently selected"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"model": model})
}

func UpdateTrainingSettings(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ID                string  `json:"id" binding:"required"`
		ExpectedAccuracy  float32 `json:"expected_accuracy" binding:"required"`
		ExpectedPrecision float32 `json:"expected_precision" binding:"required"`
		ExpectedRecall    float32 `json:"expected_recall" binding:"required"`
		ExpectedF1        float32 `json:"expected_f1" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := config.DB.Model(&models.AIModel{}).Where("id = ?", input.ID).Updates(map[string]interface{}{
		"expected_accuracy":  input.ExpectedAccuracy,
		"expected_precision": input.ExpectedPrecision,
		"expected_recall":    input.ExpectedRecall,
		"expected_f1":        input.ExpectedF1,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update training setting stats"})
		return
	}

	config.ModelSettingUpdated = true

	c.JSON(http.StatusOK, gin.H{"message": "training setting updated"})
}

func UpdateTrainingTime(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ID         string  `json:"id" binding:"required"`
		TrainEvery float64 `json:"train_every"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := config.DB.Model(&models.AIModel{}).Where("id = ?", input.ID).Updates(map[string]interface{}{
		"train_every": input.TrainEvery,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update training time stats"})
		return
	}

	config.ModelSettingUpdated = true

	c.JSON(http.StatusOK, gin.H{"message": "training time updated"})
}
