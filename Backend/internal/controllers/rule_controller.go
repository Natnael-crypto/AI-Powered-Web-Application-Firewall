package controllers

import (
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddRule(c *gin.Context) {
	if err := services.AddRuleService(c); err != nil {
		c.JSON(err.Status, gin.H{"error": err.Message})
	}
}

func GetRules(c *gin.Context) {
	data, err := services.GetRulesService(c)
	if err != nil {
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rules": data})
}

func GetAllRulesAdmin(c *gin.Context) {
	data, err := services.GetAllRulesAdminService(c)
	if err != nil {
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rules": data})
}

func GetOneRule(c *gin.Context) {
	data, err := services.GetOneRuleService(c)
	if err != nil {
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}
	c.JSON(http.StatusOK, data)
}

func UpdateRule(c *gin.Context) {
	if err := services.UpdateRuleService(c); err != nil {
		c.JSON(err.Status, gin.H{"error": err.Message})
	}
}

func ActivateRule(c *gin.Context) {
	if err := services.ToggleRuleActivation(c, true); err != nil {
		c.JSON(err.Status, gin.H{"error": err.Message})
	}
}

func DeactivateRule(c *gin.Context) {
	if err := services.ToggleRuleActivation(c, false); err != nil {
		c.JSON(err.Status, gin.H{"error": err.Message})
	}
}

func DeleteRule(c *gin.Context) {
	if err := services.DeleteRuleService(c); err != nil {
		c.JSON(err.Status, gin.H{"error": err.Message})
	}
}

func GetRuleMetadata(c *gin.Context) {
	services.GetRuleMetadataService(c)
}
