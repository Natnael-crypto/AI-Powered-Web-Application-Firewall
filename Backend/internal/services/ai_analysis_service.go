package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	analysisQueue = make(map[string]bool)
	queueMutex    sync.Mutex
)

func QueueRequestService(c *gin.Context) (gin.H, int) {
	var input struct {
		RequestID string `json:"request_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": "request_id is required"}, http.StatusBadRequest
	}

	request, err := repository.GetRequestByID(input.RequestID)
	if err != nil {
		return gin.H{"error": "request not found"}, http.StatusNotFound
	}

	app, err := repository.GetApplicationByName(request.ApplicationName)
	if err != nil {
		return gin.H{"error": "application not found"}, http.StatusNotFound
	}

	if c.GetString("role") != "super_admin" &&
		!utils.HasAccessToApplication(utils.GetAssignedApplicationIDs(c), app.ApplicationID) {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	queueMutex.Lock()
	analysisQueue[input.RequestID] = true
	queueMutex.Unlock()

	return gin.H{"message": "Request queued for AI analysis"}, http.StatusOK
}

func FetchRequestsService(c *gin.Context) (gin.H, int) {
	queueMutex.Lock()
	var ids []string
	for id := range analysisQueue {
		ids = append(ids, id)
	}
	analysisQueue = make(map[string]bool)
	queueMutex.Unlock()

	if len(ids) == 0 {
		return gin.H{"message": "No requests in queue"}, http.StatusOK
	}

	requests, err := repository.GetRequestsByIDs(ids)
	if err != nil {
		return gin.H{"error": "failed to retrieve requests"}, http.StatusInternalServerError
	}

	return gin.H{"requests": requests}, http.StatusOK
}

func SubmitAnalysisResultsService(c *gin.Context) (gin.H, int) {
	var results []struct {
		RequestID  string `json:"request_id"`
		ThreatType string `json:"threat_type"`
	}
	if err := c.ShouldBindJSON(&results); err != nil {
		return gin.H{"error": "invalid input"}, http.StatusBadRequest
	}

	for _, res := range results {
		repository.UpdateRequestAnalysisResult(res.RequestID, res.ThreatType)
	}

	return gin.H{"message": "Analysis results updated"}, http.StatusOK
}

func GetModelForMLService(c *gin.Context) (gin.H, int) {
	models, err := repository.GetAllModels()
	if err != nil {
		return gin.H{"message": "no models available"}, http.StatusOK
	}
	return gin.H{"models": models}, http.StatusOK
}

func SubmitTrainResultsService(c *gin.Context) (gin.H, int) {
	var input models.AIModel
	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": "invalid input"}, http.StatusBadRequest
	}
	err := repository.UpdateModelTrainingResult(input)
	if err != nil {
		return gin.H{"error": "failed to update model stats"}, http.StatusInternalServerError
	}
	return gin.H{"message": "model results updated"}, http.StatusOK
}

func SelectActiveModelService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}
	id := c.Param("model_id")

	if err := repository.DeselectActiveModel(); err != nil {
		return gin.H{"error": "failed to deselect existing model"}, http.StatusInternalServerError
	}
	if err := repository.SelectModel(id); err != nil {
		return gin.H{"error": "failed to select model"}, http.StatusInternalServerError
	}

	config.SelectModel = true
	return gin.H{"message": "model selected successfully"}, http.StatusOK
}

func GetSelectedModelService(c *gin.Context) (gin.H, int) {
	model, err := repository.GetSelectedModel()
	if err != nil {
		return gin.H{"message": "no model currently selected"}, http.StatusOK
	}
	config.SelectModel = false
	return gin.H{"model": model}, http.StatusOK
}

func GetAllModelsService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}
	models, err := repository.GetAllModels()
	if err != nil {
		return gin.H{"message": "no model found"}, http.StatusOK
	}
	return gin.H{"model": models}, http.StatusOK
}

func UpdateTrainingSettingsService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input struct {
		ID                string  `json:"id" binding:"required"`
		ExpectedAccuracy  float32 `json:"expected_accuracy" binding:"required"`
		ExpectedPrecision float32 `json:"expected_precision" binding:"required"`
		ExpectedRecall    float32 `json:"expected_recall" binding:"required"`
		ExpectedF1        float32 `json:"expected_f1" binding:"required"`
		TrainEvery        float64 `json:"train_every" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": "invalid input"}, http.StatusBadRequest
	}

	ms := models.AIModel{
		ID:                input.ID,
		ExpectedAccuracy:  input.ExpectedAccuracy,
		ExpectedPrecision: input.ExpectedPrecision,
		ExpectedRecall:    input.ExpectedRecall,
		ExpectedF1:        input.ExpectedF1,
		TrainEvery:        input.TrainEvery * 86400000,
	}

	err := repository.UpdateTrainingSettings(ms)
	if err != nil {
		return gin.H{"error": "failed to update training setting stats"}, http.StatusInternalServerError
	}

	config.ModelSettingUpdated = true
	return gin.H{"message": "training setting updated"}, http.StatusOK
}
