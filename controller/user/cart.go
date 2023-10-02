package user

import (
	"RmuttPlace/db"
	"RmuttPlace/dto"
	"RmuttPlace/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddCart(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var json dto.CartRequest
	var cart model.Cart
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var product model.Product
	db.Conn.Find(&product, "id =?", json.ProductID)
	if product.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"My product": "Not Found"})
		return
	}
	var count int64
	db.Conn.Model(&model.Cart{}).Where("user_id=?", uint(userId)).Group("product_id").Count(&count)
	if count > 100 {
		c.JSON(http.StatusOK, gin.H{"message": "cart max"})
		return
	}
	db.Conn.Find(&cart, "product_id = ? AND user_id = ?", json.ProductID, int(userId))
	if uint(userId) == cart.UserID && json.ProductID == cart.ProductID {
		db.Conn.Model(&cart).Where("product_id = ? AND user_id = ?", json.ProductID, int(userId)).Updates(model.Cart{UserID: uint(userId), ProductID: json.ProductID,
			Quantity: json.Quantity + cart.Quantity})
		c.JSON(http.StatusOK, gin.H{"cart ==": cart, "total": count})
		return
	} else {
		cart.ProductID = json.ProductID
		cart.UserID = uint(userId)
		cart.Quantity = json.Quantity
		db.Conn.Create(&cart)
		result := dto.CartResponse{
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
		}
		c.JSON(http.StatusOK, gin.H{"cart": result, "total": count})
		return
	}
}

func MyCart(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var carts []model.Cart
	db.Conn.Model(&model.Cart{}).Preload("Product").Find(&carts, "user_id=?", uint(userId))
	c.JSON(http.StatusOK, gin.H{"cart": carts})
}

func DeleteProductMyCart(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	id := c.Param("id")
	var cart model.Cart
	db.Conn.Find(&cart, "user_id =? and product_id =?", uint(userId), id)
	if cart.ProductID == 0 && cart.UserID == 0 {
		c.JSON(http.StatusOK, gin.H{"product": "Not found"})
		return
	}
	db.Conn.Where("user_id=? and product_id=?", uint(userId), id).Delete(&cart)
	c.JSON(http.StatusOK, gin.H{"product": "Delete success"})
}

func UpdateQuantity(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var json dto.CartRequest
	var cart model.Cart
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.Conn.Find(&cart, "product_id = ? AND user_id = ?", json.ProductID, int(userId))
	if uint(userId) == cart.UserID && json.ProductID == cart.ProductID {
		db.Conn.Model(&cart).Where("product_id = ? AND user_id = ?", json.ProductID, int(userId)).Updates(model.Cart{UserID: uint(userId), ProductID: json.ProductID,
			Quantity: json.Quantity})
		c.JSON(http.StatusOK, gin.H{"cart": cart})
		return
	}
}
