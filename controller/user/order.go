package user

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderBody struct {
	StoreID uint
	Carts   []OrderItemBody
}
type OrderItemBody struct {
	// StoreID   uint
	ProductID uint
	Quantity  int
}

func CreateOrder(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var json OrderBody
	var order model.Order
	var cart model.Cart
	// var product model.Product
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
		db.Conn.Delete(&cart, "user_id =? and product_id=?", uint(userId), product.ProductID)
	}
	// 	var quantity int
	// 	for _,q := range json.Carts{
	// 		quantity = append(quantity,model.OrderItem{
	// 			Quantity: q.Quantity,
	// 		})
	// 	}
	// product.Available=product.Available-quantity
	order.StoreID = json.StoreID
	order.UserID = uint(userId)
	order.Products = orderItems
	order.ShipmentID = 0
	if err := db.Conn.Preload("Product").Create(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": order})
}

func MyOrderAll(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var order []model.Order
	var orderItem []model.OrderItem
	query := db.Conn.Find(&order, "user_id=?", uint(userId))
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
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

// type ReadOrderResponse struct {
// 	OrederID  uint
// 	ProductID uint
// 	Quantity  int
// 	Product   Product
// }

func MyOrderFindOne(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	id := c.Param("id")
	var order model.Order
	var product model.Product
	var orderItems []model.OrderItem
	query := db.Conn.Find(&order, "user_id=? and id=?", uint(userId), id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query2 := db.Conn.Find(&orderItems, "order_id=?", order.ID)
	if err := query2.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	//หา product_id เอามา query หา Product
	query3 := db.Conn.Preload("Product").Find(&product, "id=?")
	// for _, product := range orderItems {
	// 	orderItems = append(orderItems, model.OrderItem{
	// 		OrderID:   product.OrderID,
	// 		ProductID: product.ProductID,
	// 		Product: model.Product{
	// 			StoreId: product.Product.StoreId,
	// 			Name:    product.Product.Name,
	// 		},
	// 	})
	// }
	c.JSON(http.StatusOK, query3)
}
