package store

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderUpdateBody struct {
	Tracking   string `json:"tracking"`
	ShipmentID int    `json:"shipmentId"`
}
type AddTrackingOrderRead struct {
	ID           uint
	UserID       uint
	StoreID      uint
	ShipmentID   uint
	ShipmentName string
	Tracking     string
}

func AddTrackingOrder(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var order model.Order
	var shipment model.Shipment
	var json OrderUpdateBody
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Preload("Shipment").Find(&order, "store_id = ? AND id=?", uint(storeId), id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&shipment, "id=?", json.ShipmentID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	order.Tracking = json.Tracking
	order.ShipmentID = uint(json.ShipmentID)
	db.Conn.Save(&order)
	result := AddTrackingOrderRead{
		ID:         order.ID,
		UserID:     order.UserID,
		StoreID:    order.StoreID,
		ShipmentID: uint(order.ShipmentID),
		Tracking:   order.Tracking,
	}

	c.JSON(http.StatusOK, result)

}

func GetOrderAll(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var order []model.Order
	if err := db.Conn.Find(&order, "store_id = ? ", storeId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
func GetOrderOne(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var order []model.Order
	var orderItem []model.OrderItem
	if err := db.Conn.Preload("Shipment").Find(&order, "store_id = ? and id=?", storeId, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	for _, orders := range order {
		var ot []model.OrderItem
		query := db.Conn.Preload("Product").Find(&ot, "order_id=?", orders.ID)
		if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		orderItem = append(orderItem, ot...)
	}

	c.JSON(http.StatusOK, orderItem)

}
