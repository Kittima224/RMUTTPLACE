package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddTrackingOrderRead struct {
	ID           uint
	UserID       uint
	StoreID      uint
	ShipmentID   uint
	ShipmentName string
	Tracking     string
}

func GetOrderAll(c *gin.Context) {
	var orders []model.Order
	if err := db.Conn.Preload("Shipment").Find(&orders).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var result []AddTrackingOrderRead
	for _, order := range orders {
		result = append(result, AddTrackingOrderRead{
			ID:           order.ID,
			UserID:       order.UserID,
			StoreID:      order.StoreID,
			ShipmentID:   uint(order.ShipmentID),
			ShipmentName: order.Shipment.Name,
			Tracking:     order.Tracking,
		})
	}
	c.JSON(http.StatusOK, result)
}

func GetOrderOne(c *gin.Context) {
	id := c.Param("id")
	var order []model.Order
	var orderItem []model.OrderItem
	if err := db.Conn.Find(&order, " id=?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
