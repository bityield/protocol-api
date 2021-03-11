package controllers

import (
	"net/http"

	"github.com/bityield/bityield-api/infra/database/models"
	"github.com/gin-gonic/gin"
)

// var (
// 	db *gorm.DB
// )

// func init() {
// 	db = database.ConnectDataBase()
// }

// CreateFundInput validator for POST
type CreateFundInput struct {
	Name string `json:"name" binding:"required"`
}

// UpdateFundInput validator for PUT
type UpdateFundInput struct {
	Name string `json:"name" binding:"required"`
}

// FindFunds returns all funds
func FindFunds(c *gin.Context) {
	db, err := GetConn(c)
	if err != nil {
		panic(err)
	}

	var funds []models.Fund
	db.Find(&funds)

	c.JSON(http.StatusOK, gin.H{"data": funds})
}

// CreateFund post creation
func CreateFund(c *gin.Context) {
	db, err := GetConn(c)
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

// FindFund returns a sinlg fund by Id
func FindFund(c *gin.Context) {
	db, err := GetConn(c)
	if err != nil {
		panic(err)
	}

	var fund models.Fund

	if err := db.Where("id = ?", c.Param("id")).First(&fund).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fund})
}

// UpdateFund ...
func UpdateFund(c *gin.Context) {
	db, err := GetConn(c)
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
