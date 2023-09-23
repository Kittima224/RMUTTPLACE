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
	Tracking string `json:"tracking"`
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var order model.Order

	var json OrderUpdateBody
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&order, "store_id = ? and id=?", storeId, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	order.Tracking = json.Tracking
	db.Conn.Save(&order)
	c.JSON(http.StatusOK, order)
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
	if err := db.Conn.Find(&order, "store_id = ? and id=?", storeId, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
