package user

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Favorite struct {
	UserID    uint `json:"userId"`
	ProductID uint `json:"productId"`
}

func AddUnFav(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var json Favorite
	var fav model.Favorite
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
	db.Conn.First(&fav, "product_id=? AND user_id=?", json.ProductID, uint(userId))
	if uint(userId) == fav.UserID && fav.ProductID == json.ProductID {
		db.Conn.Where("user_id=? and product_id=?", int(userId), json.ProductID).Delete(&fav)
		c.JSON(http.StatusOK, gin.H{"message": "Un Favorite"})
		return
	}
	fav.ProductID = json.ProductID
	fav.UserID = uint(userId)
	db.Conn.Create(&fav)
	c.JSON(http.StatusOK, gin.H{"cart": fav, "message": "success"})
}

func MyFav(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var favs []model.Favorite
	db.Conn.Model(&Favorite{}).Preload("Product").Find(&favs, "user_id=?", uint(userId))
	c.JSON(http.StatusOK, gin.H{"cart": favs})
}
