package user

import (
	"RmuttPlace/db"
	"RmuttPlace/dto"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var json dto.OrderRequest
	//var product model.Product
	var user model.User
	var order model.Order

	if err := db.Conn.First(&user, userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderItems []model.OrderItem
	for _, product := range json.Carts {
		orderItems = append(orderItems, model.OrderItem{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		})
	}
	order.StoreID = json.StoreID
	order.UserID = uint(userId)
	order.Products = orderItems
	order.ShipmentID = 3

	if err := db.Conn.Create(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := dto.OrderRead{
		OrderID:    order.ID,
		ShipmentID: order.ShipmentID,
	}
	db.Conn.Find(&orderItems, "order_id=?", order.ID)

	var Available []struct {
		ID        uint
		Available int
	}
	var product model.Product

	db.Conn.Raw("SELECT  products.id as id,(products.available-order_items.quantity) as available FROM orders JOIN order_items on order_items.order_id=orders.id JOIN products on order_items.product_id=products.id WHERE orders.id=?", order.ID).Scan(&Available)

	for _, item := range Available {
		db.Conn.Raw("UPDATE PRODUCTS SET AVAILABLE =? WHERE ID=?", item.Available, item.ID).Scan(&product)
	}
	c.JSON(http.StatusOK, gin.H{"order": result})
}

func MyOrderAll(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var orders []model.Order

	var user model.User
	if err := db.Conn.First(&user, userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query := db.Conn.Preload("Store").Preload("Shipment").Find(&orders, "user_id=?", uint(userId))
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
			Store: dto.StoreRead{
				ID:   order.StoreID,
				Name: order.Store.NameStore,
			},
		})
	}

	c.JSON(http.StatusOK, result)

}

func MyOrderFindOne(c *gin.Context) {
	userId := c.MustGet("userId").(float64)

	var user model.User
	if err := db.Conn.First(&user, userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	var order model.Order
	var orderItems []model.OrderItem
	query := db.Conn.Preload("Shipment").Preload("Store").Find(&order, "user_id=? and id=?", uint(userId), id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query2 := db.Conn.Preload("Product").Find(&orderItems, "order_id=?", order.ID)
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
