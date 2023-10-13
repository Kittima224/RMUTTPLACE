package store

import (
	"RmuttPlace/db"
	"RmuttPlace/dto"
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
type OrderUpdateBody struct {
	Tracking   string `json:"tracking"`
	ShipmentID uint   `json:"shipmentId"`
}

func AddTrackingOrder(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var order model.Order
	var json OrderUpdateBody
	var orderItems []model.OrderItem
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Preload("Store").Preload("Shipment").Find(&order, "store_id = ? AND id=?", uint(storeId), id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query2 := db.Conn.Preload("Product").Find(&orderItems, "order_id", id)
	if err := query2.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	order.Tracking = json.Tracking
	order.ShipmentID = json.ShipmentID
	db.Conn.Save(&order)
	db.Conn.Model(&order).Updates(OrderUpdateBody{
		Tracking:   json.Tracking,
		ShipmentID: json.ShipmentID,
	})
	result := dto.OrderReadOne{
		ID: order.ID,
		Store: dto.StoreRead{
			ID:   order.Store.ID,
			Name: order.Store.NameStore,
		},
		Shipment: dto.ShipmentRead{
			ID:   order.ShipmentID,
			Name: order.Shipment.Name,
		},
		Tracking: order.Tracking,
	}
	var ot []dto.OrderItemRead
	for _, o := range orderItems {
		ot = append(ot, dto.OrderItemRead{
			ID:       o.Product.ID,
			Image:    o.Product.Image,
			Price:    o.Product.Price,
			Quantity: o.Quantity,
			Name:     o.Product.Name,
		})
	}
	result.Products = ot
	c.JSON(http.StatusOK, result)

}

func GetOrderAll(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var orders []model.Order
	var store model.Store
	if err := db.Conn.Find(&store, storeId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Preload("Shipment").Find(&orders, "store_id = ? ", storeId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var result []dto.OrderReadAll
	for _, order := range orders {
		result = append(result, dto.OrderReadAll{
			ID:           order.ID,
			UserID:       order.UserID,
			ShipmentID:   uint(order.ShipmentID),
			ShipmentName: order.Shipment.Name,
			Tracking:     order.Tracking,
		})
	}
	c.JSON(http.StatusOK, result)
}
func GetOrderOne(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var store model.Store
	if err := db.Conn.Find(&store, storeId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var order model.Order
	var orderItems []model.OrderItem
	if err := db.Conn.Preload("User").Preload("Store").Preload("Shipment").Find(&order, "store_id = ? and id=?", storeId, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query2 := db.Conn.Preload("Product").Find(&orderItems, "order_id=?", id)
	if err := query2.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	result := dto.OrderReadOne{
		ID: order.ID,
		Store: dto.StoreRead{
			ID:   order.Store.ID,
			Name: order.Store.NameStore,
		},
		Shipment: dto.ShipmentRead{
			ID:   order.ShipmentID,
			Name: order.Shipment.Name,
		},
		User: dto.UserAddress{
			ID:          order.User.ID,
			Name:        order.User.UserName,
			Image:       order.User.Image,
			Address:     order.User.Address,
			District:    order.User.District,
			SubDistrict: order.User.SubDistrict,
			Province:    order.User.Province,
			Zipcode:     order.User.Zipcode,
			Tel:         order.User.Tel,
		},
		Tracking: order.Tracking,
	}
	var ot []dto.OrderItemRead
	for _, o := range orderItems {
		ot = append(ot, dto.OrderItemRead{
			ID:       o.Product.ID,
			Image:    o.Product.Image,
			Price:    o.Product.Price,
			Quantity: o.Quantity,
			Name:     o.Product.Name,
		})
	}
	result.Products = ot
	c.JSON(http.StatusOK, result)

}
