package models

import "time"

type AIModel struct {
	ID                    string    `json:"id" gorm:"primaryKey"`
	NumberRequestsUsed    int       `json:"number_requests_used" gorm:"not null"`
	ModelsName            string    `json:"models_name" gorm:"not null;unique"`
	PercentTrainData      float32   `json:"percent_train_data" gorm:"not null"`
	PercentNormalRequests float32   `json:"percent_normal_requests" gorm:"not null"`
	NumTrees              int       `json:"num_trees" gorm:"not null"`         // Number of trees (n_estimators)
	MaxDepth              int       `json:"max_depth" gorm:"not null"`         // Max depth of each tree
	MinSamplesSplit       int       `json:"min_samples_split" gorm:"not null"` // Min samples to split a node
	MinSamplesLeaf        int       `json:"min_samples_leaf" gorm:"not null"`  // Min samples at leaf node
	MaxFeatures           string    `json:"max_features" gorm:"not null"`      // Features to consider at each split (e.g., "sqrt")
	Criterion             string    `json:"criterion" gorm:"not null"`         // Split criterion (e.g., "gini", "entropy")
	Accuracy              float32   `json:"accuracy"`
	Precision             float32   `json:"precision"`
	Recall                float32   `json:"recall"`
	F1                    float32   `json:"f1"`
	Selected              bool      `json:"selected" gorm:"not null"`
	Modeled               bool      `json:"modeled" gorm:"not null"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
