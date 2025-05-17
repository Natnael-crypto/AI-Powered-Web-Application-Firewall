package models

import "time"

type AIModel struct {
	ID                    string    `json:"id" gorm:"primaryKey"`
	ModelsName            string    `json:"models_name" gorm:"not null;unique"`
	Accuracy              float32   `json:"accuracy"`
	Precision             float32   `json:"precision"`
	Recall                float32   `json:"recall"`
	F1                    float32   `json:"f1"`
	ExpectedAccuracy      float32   `json:"expected_accuracy"`
	ExpectedPrecision     float32   `json:"expected_precision"`
	ExpectedRecall        float32   `json:"expected_recall"`
	ExpectedF1            float32   `json:"expected_f1"`
	Selected              bool      `json:"selected" gorm:"not null"`
	UpdatedAt             time.Time `json:"updated_at"`
	LastTrainedAt         time.Time `json:"last_trained_at"`
	TrainEvery            float64   `json:"train_every"`
	PercentOfTrainingData float64   `json:"percent_of_training_data"`
	ModelType             string    `json:"model_type"`
}
