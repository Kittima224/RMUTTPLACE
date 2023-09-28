package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderReadAll struct {
	ID           uint
	UserID       uint
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
	var result []OrderReadAll
	for _, order := range orders {
		result = append(result, OrderReadAll{
			ID:           order.ID,
			UserID:       order.UserID,
			ShipmentID:   uint(order.ShipmentID),
			ShipmentName: order.Shipment.Name,
			Tracking:     order.Tracking,
		})
	}
	c.JSON(http.StatusOK, result)
}

type OrderReadOne struct {
	ID       uint
	Store    model.StoreRead
	Products []OrderItemRead
}

// type OrderItemRead struct {
// 	Image    string
// 	Name     string
// 	Quantity int
// 	Price    int
// }

type OrderItemRead struct {
	ProductID uint
	Image     string
	Price     int
	Quantity  int
}

func GetOrderOne(c *gin.Context) {
	id := c.Param("id")
	var order model.Order
	var orderItems []model.OrderItem
	query := db.Conn.Preload("Store").Find(&order, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query2 := db.Conn.Preload("Product").Find(&orderItems, "order_id", id)
	if err := query2.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	result := OrderReadOne{
		ID: order.ID,
		Store: model.StoreRead{
			ID:   order.Store.ID,
			Name: order.Store.NameStore,
		},
	}
	var ot []OrderItemRead
	for _, o := range orderItems {
		ot = append(ot, OrderItemRead{
			ProductID: o.Product.ID,
			Image:     o.Product.Image,
			Price:     o.Product.Price,
			Quantity:  o.Quantity,
		})
	}
	result.Products = ot
	c.JSON(http.StatusOK, result)
}
