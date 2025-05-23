package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"log"
	"time"
)

func CreateModelTrainingRequest(model models.AIModel) error {
	if err := config.DB.Create(&model).Error; err != nil {
		return err
	}
	return nil
}

func CreateModel() {

	models_name := []string{"Radome Forest", "Support Vector Machines", "XGBoost"}

	for _, model := range models_name {
		AiModel := models.AIModel{
			ID:                    GenerateUUID(),
			ModelsName:            model,
			Accuracy:              0.0,
			Precision:             0.0,
			Recall:                0.0,
			F1:                    0.0,
			ExpectedAccuracy:      0.0,
			ExpectedPrecision:     0.0,
			ExpectedRecall:        0.0,
			ExpectedF1:            0.0,
			Selected:              false,
			UpdatedAt:             time.Now(),
			LastTrainedAt:         time.Now(),
			TrainEvery:            86400000.0,
			PercentOfTrainingData: 0.0,
			ModelType:             model,
		}
		if err := CreateModelTrainingRequest(AiModel); err != nil {
			log.Println(err)
		}
	}
}
