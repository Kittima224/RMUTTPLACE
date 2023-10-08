package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShipmentBody struct {
	Name string `json:"name" binding:"required"`
}

func ShipmentCreate(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var json ShipmentBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shipment := model.Shipment{
		Name: json.Name,
	}
	db.Conn.Find(&shipment, "name =?", json.Name)
	if shipment.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Duplicate",
		})
		return
	}
	if err := db.Conn.Create(&shipment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"shipment": shipment})
}

func ShipmentAll(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var shipment []model.Shipment
	db.Conn.Find(&shipment)
	c.JSON(http.StatusOK, shipment)
}

func ShipmentOne(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	var shipment model.Shipment
	if err := db.Conn.First(&shipment, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "Shipment Read Success", "shipment": shipment})
}

func ShipmentUpdate(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	var shipment model.Shipment
	var json ShipmentBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.First(&shipment, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	shipment.Name = json.Name
	db.Conn.Save(&shipment)
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "update shipment", "shipment": shipment})
}

type shipmentid struct {
	Id uint `json:"id"`
}

func ShipmentDel(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var shipment model.Shipment
	var json shipmentid
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.Conn.Find(&shipment, "id =?", json.Id)
	if shipment.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Not Found"})
		return
	} else {
		db.Conn.Delete(&shipment).Where("id =?", json.Id)
		c.JSON(http.StatusOK, gin.H{"message": "Delete user ?"})
		return
	}
}
