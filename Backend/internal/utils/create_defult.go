package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateModelTrainingRequest(model models.AIModel) error {
	if err := config.DB.Create(&model).Error; err != nil {
		return err
	}
	return nil
}

func CreateModel() {

	models_name := []string{"Radome Forest"}

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

func AddNotificationRule() {
	defaultRules := []struct {
		Name       string
		ThreatType string
	}{
		{"SQL Injection Attack", "sql"},
		{"XSS Attack", "xss"},
		{"Command Injection Attack", "cmd"},
		{"LDAP Injection Attack", "ldap"},
		{"NoSQL Injection Attack", "nosql"},
		{"Code Injection Attack", "code"},
		{"Path Traversal Attack", "path"},
		{"Rate Limited", "rate"},
	}

	threshold := 10
	timeWindow := 10
	isActive := true

	for _, def := range defaultRules {
		var existing models.NotificationRule
		err := config.DB.Where("threat_type = ?", def.ThreatType).First(&existing).Error

		if err == nil {
			continue
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("Error checking existing rule: %v\n", err)
			continue
		}

		rule := models.NotificationRule{
			ID:         uuid.New().String(),
			Name:       def.Name,
			ThreatType: def.ThreatType,
			Threshold:  threshold,
			TimeWindow: timeWindow,
			IsActive:   isActive,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := config.DB.Create(&rule).Error; err != nil {
			fmt.Printf("Error creating rule %s: %v\n", def.Name, err)
		}
	}
}
