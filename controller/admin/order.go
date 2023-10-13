package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/dto"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrderAll(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var orders []model.Order
	if err := db.Conn.Preload("Shipment").Find(&orders).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	var order model.Order
	var orderItems []model.OrderItem
	query := db.Conn.Preload("User").Preload("Store").Find(&order, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query2 := db.Conn.Preload("Product").Find(&orderItems, "order_id", id)
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
