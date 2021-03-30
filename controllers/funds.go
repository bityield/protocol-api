package controllers

import (
	"net/http"

	"github.com/bityield/protocol-api/infra/database/models"
	"github.com/gin-gonic/gin"
)

// CreateFundInput validator for POST
type CreateFundInput struct {
	Name string `json:"name" binding:"required"`
}

// UpdateFundInput validator for PUT
type UpdateFundInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateFund post creation
func CreateFund(c *gin.Context) {
	db, err := GetDatabase(c)
	if err != nil {
		panic(err)
	}

	var input CreateFundInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fund := models.Fund{Name: input.Name}
	db.Create(&fund)

	c.JSON(http.StatusOK, gin.H{"data": fund})
}

// FindFunds returns all funds
func FindFunds(c *gin.Context) {
	db, err := GetDatabase(c)
	if err != nil {
		panic(err)
	}

	var funds []models.Fund
	db.Preload("Assets").Find(&funds)

	c.JSON(http.StatusOK, gin.H{"funds": funds})
}

// FindFund returns a sinlg fund by Id
func FindFund(c *gin.Context) {
	db, err := GetDatabase(c)
	if err != nil {
		panic(err)
	}

	var fund models.Fund

	if err := db.Where("address = ?", c.Param("id")).Preload("Assets").First(&fund).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fund": fund})
}

// UpdateFund ...
func UpdateFund(c *gin.Context) {
	db, err := GetDatabase(c)
	if err != nil {
		panic(err)
	}

	var fund models.Fund
	if err := db.Where("id = ?", c.Param("id")).First(&fund).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateFundInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&fund).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": fund})
}
