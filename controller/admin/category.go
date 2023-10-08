package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryBody struct {
	Name string `json:"name" binding:"required"`
}

func Create(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	var json CategoryBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	category := model.Category{
		Name: json.Name,
	}
	db.Conn.Find(&category, "name =?", json.Name)
	if category.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Duplicate",
		})
		return
	}
	if err := db.Conn.Create(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "create category", "category": category})
}

func CategoryAll(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	var categories []model.Category
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&categories).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}
func CategoryOne(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	id := c.Param("id")
	var category model.Category
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "Category Read Success", "category": category})
}
func CategoryUpdate(c *gin.Context) {
	id := c.Param("id")
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	var category model.Category
	var json CategoryBody
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	category.Name = json.Name
	db.Conn.Save(&category)
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "update category", "category": category})
}

type Categoryid struct {
	Id uint `json:"id"`
}

func CategoryDel(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	var category model.Category
	var json Categoryid
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	db.Conn.Find(&category, "id =?", json.Id)
	if category.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Not Found"})
		return
	} else {
		db.Conn.Delete(&category).Where("id =?", json.Id)
		c.JSON(http.StatusOK, gin.H{"message": "Delete user ?"})
		return
	}
}
