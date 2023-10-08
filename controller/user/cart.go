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
	db.Conn.Find(&cart, "product_id = ? AND user_id = ?", json.ProductID, userId)
	if uint(userId) == cart.UserID && json.ProductID == cart.ProductID {
		db.Conn.Model(&cart).Where("product_id = ? AND user_id = ?", json.ProductID, userId)
		cart.Quantity = cart.Quantity + json.Quantity
		db.Conn.Save(&cart)
		result := dto.CartResponse{
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
		}
		c.JSON(http.StatusOK, gin.H{"carts": result, "total": count})
		return
	} else {
		cart.ProductID = json.ProductID
		cart.StoreID = json.StoreID
		cart.UserID = uint(userId)
		cart.Quantity = json.Quantity
		db.Conn.Create(&cart)
		result := dto.CartResponse{
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
		}
		c.JSON(http.StatusOK, gin.H{"carts": result, "total": count + 1})
		return
	}
}

func MyCart(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var carts []model.Cart
	var cart model.Cart
	var store model.Store
	db.Conn.Model(&model.Cart{}).Preload("Store").Preload("Product").Find(&carts, "user_id=?", uint(userId))

	if err := db.Conn.Find(&store, "id =?", cart.StoreID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var result []dto.ReadProductInCart
	for _, p := range carts {
		result = append(result, dto.ReadProductInCart{
			Store: dto.StoreRead{
				ID:   p.StoreID,
				Name: p.Store.NameStore,
			},
			ID:       p.ProductID,
			Image:    p.Product.Image,
			Name:     p.Product.Name,
			Quantity: p.Quantity,
		})
	}
	c.JSON(http.StatusOK, gin.H{"carts": result})
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
