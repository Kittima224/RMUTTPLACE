package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrderAll(c *gin.Context) {
	var order []model.Order
	if err := db.Conn.Find(&order).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
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
