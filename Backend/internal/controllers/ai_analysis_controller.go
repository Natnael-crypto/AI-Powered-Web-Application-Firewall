package controllers

import (
	"net/http"
	"sync"
	"time"

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

func CreateModelTrainingRequest(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}
	var input struct {
		ModelsName            string  `json:"models_name" binding:"required"`
		NumberRequestsUsed    int     `json:"number_requests_used" binding:"required"`
		PercentTrainData      float32 `json:"percent_train_data" binding:"required"`
		PercentNormalRequests float32 `json:"percent_normal_requests" gorm:"not null"`
		NumTrees              int     `json:"num_trees" binding:"required"`
		MaxDepth              int     `json:"max_depth" binding:"required"`
		MaxFeatures           string  `json:"max_features" binding:"required"`
		MinSamplesSplit       int     `json:"min_samples_split" binding:"required"`
		MinSamplesLeaf        int     `json:"min_samples_leaf" binding:"required"`
		Criterion             string  `json:"criterion" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	ai_model := models.AIModel{
		ID:                    utils.GenerateUUID(),
		ModelsName:            input.ModelsName,
		NumberRequestsUsed:    input.NumberRequestsUsed,
		PercentTrainData:      input.PercentTrainData,
		PercentNormalRequests: input.PercentNormalRequests,
		NumTrees:              input.NumTrees,
		MaxDepth:              input.MaxDepth,
		MinSamplesSplit:       input.MinSamplesSplit,
		MaxFeatures:           input.MaxFeatures,
		MinSamplesLeaf:        input.MinSamplesLeaf,
		Criterion:             input.Criterion,
		Accuracy:              0,
		Precision:             0,
		Recall:                0,
		F1:                    0,
		Selected:              false,
		Modeled:               false,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if err := config.DB.Create(&ai_model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create model"})
		return
	}

	config.UntrainedModel = true

	c.JSON(http.StatusOK, gin.H{"message": "model training request created", "model": ai_model})
}

func GetUntrainedModelForML(c *gin.Context) {
	var model models.AIModel

	if err := config.DB.Where("modeled = ?", false).First(&model).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "no untrained model available"})
		return
	}

	totalRequests := model.NumberRequestsUsed
	normalCount := int(float32(totalRequests) * model.PercentNormalRequests / 100)
	maliciousCount := totalRequests - normalCount

	type RequestData struct {
		ApplicationName string `json:"application_name"`
		RequestMethod   string `json:"request_method"`
		RequestURL      string `json:"request_url"`
		Headers         string `json:"headers"`
		Body            string `json:"body"`
	}

	var normalRequests []RequestData
	if err := config.DB.
		Model(&models.Request{}).
		Select("application_name, request_method, request_url, headers, body").
		Where("threat_detected = ?", false).
		Limit(normalCount).
		Scan(&normalRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch normal requests"})
		return
	}

	var maliciousRequests []RequestData
	if err := config.DB.
		Model(&models.Request{}).
		Select("application_name, request_method, request_url, headers, body").
		Where("threat_detected = ?", true).
		Limit(maliciousCount).
		Scan(&maliciousRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch malicious requests"})
		return
	}

	config.UntrainedModel = false
	c.JSON(http.StatusOK, gin.H{
		"id":                      model.ID,
		"name":                    model.ModelsName,
		"number_requests_used":    model.NumberRequestsUsed,
		"percent_train_data":      model.PercentTrainData,
		"percent_normal_requests": model.PercentNormalRequests,
		"num_trees":               model.NumTrees,
		"max_depth":               model.MaxDepth,
		"min_samples_split":       model.MinSamplesSplit,
		"min_samples_leaf":        model.MinSamplesLeaf,
		"max_features":            model.MaxFeatures,
		"criterion":               model.Criterion,
		"normal_requests":         normalRequests,
		"malicious_request":       maliciousRequests,
	})
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

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	modelID := c.Param("model_id")

	if err := config.DB.Where("modeled = ? AND id = ?", true, modelID).First(&models.AIModel{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "model not found or not modeled"})
	}

	if err := config.DB.Model(&models.AIModel{}).Where("selected = ?", true).Update("selected", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deselect existing model"})
		return
	}

	if err := config.DB.Model(&models.AIModel{}).Where("id = ?", modelID).Update("selected", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to select model"})
		return
	}

	config.SelecteModel = true

	c.JSON(http.StatusOK, gin.H{"message": "model selected successfully"})
}

func GetSelectedModel(c *gin.Context) {
	var model models.AIModel

	if err := config.DB.Where("selected = ?", true).First(&model).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "no model currently selected"})
		return
	}

	config.SelecteModel = false

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

func DeleteModel(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}
	modelsID := c.Param("model_id")
	if err := config.DB.Where("id = ?", modelsID).Delete(&models.AIModel{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "model deleted successfully",
	})

}
