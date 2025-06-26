package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

func GetRequestByID(id string) (models.Request, error) {
	var r models.Request
	err := config.DB.Where("request_id = ?", id).First(&r).Error
	return r, err
}

// func GetApplicationByName(name string) (models.Application, error) {
// 	var a models.Application
// 	err := config.DB.Where("request_id = ?", name).First(&a).Error
// 	return a, err
// }

func GetRequestsByIDs(ids []string) ([]models.Request, error) {
	var requests []models.Request
	err := config.DB.Where("request_id IN ?", ids).Find(&requests).Error
	return requests, err
}

func UpdateRequestAnalysisResult(id, threatType string) {
	config.DB.Model(&models.Request{}).Where("request_id = ?", id).Updates(map[string]interface{}{
		"ai_result":      true,
		"ai_threat_type": threatType,
	})
}

func GetAllModels() ([]models.AIModel, error) {
	var models []models.AIModel
	err := config.DB.Find(&models).Error
	return models, err
}

func UpdateModelTrainingResult(m models.AIModel) error {
	return config.DB.Model(&models.AIModel{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"accuracy":   m.Accuracy,
		"precision":  m.Precision,
		"recall":     m.Recall,
		"f1":         m.F1,
		"model_type": m.ModelType,
	}).Error
}

func DeselectActiveModel() error {
	return config.DB.Model(&models.AIModel{}).Where("selected = ?", true).Update("selected", false).Error
}

func SelectModel(id string) error {
	return config.DB.Model(&models.AIModel{}).Where("id = ?", id).Update("selected", true).Error
}

func GetSelectedModel() (models.AIModel, error) {
	var model models.AIModel
	err := config.DB.Where("selected = ?", true).First(&model).Error
	return model, err
}

func UpdateTrainingSettings(m models.AIModel) error {
	return config.DB.Model(&models.AIModel{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"expected_accuracy":  m.ExpectedAccuracy,
		"expected_precision": m.ExpectedPrecision,
		"expected_recall":    m.ExpectedRecall,
		"expected_f1":        m.ExpectedF1,
		"train_every":        m.TrainEvery,
	}).Error
}
