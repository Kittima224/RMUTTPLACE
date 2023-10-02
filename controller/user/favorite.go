package user

import (
	"RmuttPlace/db"
	"RmuttPlace/dto"
	"RmuttPlace/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddFav(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var json dto.FavoriteRequest
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
	fav.ProductID = json.ProductID
	fav.UserID = uint(userId)
	db.Conn.Create(&fav)
	result := dto.FavoriteReponse{
		ProductID: fav.ProductID,
	}
	c.JSON(http.StatusOK, gin.H{"fav": result, "message": "success"})
}

func UnFav(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var json dto.FavoriteRequest
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
	db.Conn.Delete(&fav).Where("user_id =? and product_id=?", uint(userId), json.ProductID)
	result := dto.FavoriteReponse{
		ProductID: fav.ProductID,
	}
	c.JSON(http.StatusOK, gin.H{"fav": result, "message": "delete success"})
}

func MyFav(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var favs []model.Favorite
	db.Conn.Model(&model.Favorite{}).Preload("Product").Find(&favs, "user_id=?", uint(userId))
	c.JSON(http.StatusOK, gin.H{"fav": favs})
}
