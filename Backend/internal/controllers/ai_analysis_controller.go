package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func QueueRequestForAnalysis(c *gin.Context) {
	resp, status := services.QueueRequestService(c)
	c.JSON(status, resp)
}

func FetchAndAnalyzeRequests(c *gin.Context) {
	resp, status := services.FetchRequestsService(c)
	c.JSON(status, resp)
}

func SubmitAnalysisResults(c *gin.Context) {
	resp, status := services.SubmitAnalysisResultsService(c)
	c.JSON(status, resp)
}

func GetModelForMLs(c *gin.Context) {
	resp, status := services.GetModelForMLService(c)
	c.JSON(status, resp)
}

func SubmitTrainResults(c *gin.Context) {
	resp, status := services.SubmitTrainResultsService(c)
	c.JSON(status, resp)
}

func SelectActiveModel(c *gin.Context) {
	resp, status := services.SelectActiveModelService(c)
	c.JSON(status, resp)
}

func GetSelectedModel(c *gin.Context) {
	resp, status := services.GetSelectedModelService(c)
	c.JSON(status, resp)
}

func GetModels(c *gin.Context) {
	resp, status := services.GetAllModelsService(c)
	c.JSON(status, resp)
}

func UpdateTrainingSettings(c *gin.Context) {
	resp, status := services.UpdateTrainingSettingsService(c)
	c.JSON(status, resp)
}
