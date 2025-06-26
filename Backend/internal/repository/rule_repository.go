package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

func CreateRule(rule *models.Rule, appIDs []string) error {
	if err := config.DB.Create(rule).Error; err != nil {
		return err
	}
	return AddRuleToApps(rule.RuleID, appIDs)
}

func GetRuleIDsByApp(appID string) ([]string, error) {
	var mappings []models.RuleToApp
	if err := config.DB.Where("application_id = ?", appID).Find(&mappings).Error; err != nil {
		return nil, err
	}
	var ids []string
	for _, m := range mappings {
		ids = append(ids, m.RuleID)
	}
	return ids, nil
}

func GetActiveRulesByIDs(ids []string) ([]models.Rule, error) {
	var rules []models.Rule
	if err := config.DB.Where("rule_id IN ? AND is_active = true", ids).Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func GetRuleToAppsByAppIDs(appIDs []string) ([]models.RuleToApp, error) {
	var mappings []models.RuleToApp
	if err := config.DB.Where("application_id IN ?", appIDs).Find(&mappings).Error; err != nil {
		return nil, err
	}
	return mappings, nil
}

func GetRulesByIDs(ids []string) ([]models.Rule, error) {
	var rules []models.Rule
	if err := config.DB.Where("rule_id IN ?", ids).Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func GetAppsByIDs(appIDs []string) ([]models.Application, error) {
	var apps []models.Application
	if err := config.DB.Where("application_id IN ?", appIDs).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func GetRuleByIDAndUser(ruleID, userID string) (models.Rule, error) {
	var rule models.Rule
	err := config.DB.Where("rule_id = ? AND created_by = ?", ruleID, userID).First(&rule).Error
	return rule, err
}

func DeleteRuleToApps(ruleID string) error {
	return config.DB.Where("rule_id = ?", ruleID).Delete(&models.RuleToApp{}).Error
}

func AddRuleToApps(ruleID string, appIDs []string) error {
	for _, id := range appIDs {
		m := models.RuleToApp{RuleID: ruleID, ApplicationID: id}
		if err := config.DB.Create(&m).Error; err != nil {
			return err
		}
	}
	return nil
}

func SaveRule(rule *models.Rule) error {
	return config.DB.Save(rule).Error
}

func GetAppsByRuleID(ruleID string) ([]models.Application, error) {
	var mappings []models.RuleToApp
	if err := config.DB.Where("rule_id = ?", ruleID).Find(&mappings).Error; err != nil {
		return nil, err
	}
	appIDs := make([]string, 0, len(mappings))
	for _, m := range mappings {
		appIDs = append(appIDs, m.ApplicationID)
	}
	return GetAppsByIDs(appIDs)
}

func DeleteRule(ruleID, userID string) error {
	return config.DB.Where("rule_id = ? AND created_by = ?", ruleID, userID).Delete(&models.Rule{}).Error
}
