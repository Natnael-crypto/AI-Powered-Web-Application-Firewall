package models

import "time"

type AIModel struct {
	ID                 string    `json:"id" gorm:"primaryKey"`
	NumberRequestsUsed int       `json:"number_requests_used" gorm:"not null"`
	PercentTrainData   float32   `json:"percent_train_data" gorm:"not null"`
	NumTrees           int       `json:"num_trees" gorm:"not null"`
	MaxDepth           int       `json:"max_depth" gorm:"not null"`
	MinSamplesSplit    int       `json:"min_samples_split" gorm:"not null"`
	Criterion          string    `json:"criterion" gorm:"not null"`
	Accuracy           float32   `json:"accuracy"`
	Precision          float32   `json:"precision"`
	Recall             float32   `json:"recall"`
	F1                 float32   `json:"f1"`
	Selected           bool      `json:"selected" gorm:"not null"`
	Modeled            bool      `json:"modeled" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
