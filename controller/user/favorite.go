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

	db.Conn.Find(&fav, "user_id=? and product_id=?", userId, json.ProductID)
	if fav.ID == 0 {
		fav.ProductID = json.ProductID
		fav.UserID = uint(userId)
		db.Conn.Create(&fav)
		result := dto.FavoriteReponse{
			ProductID: fav.ProductID,
		}
		c.JSON(http.StatusOK, gin.H{"fav": result, "message": "favorite product"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "duplicate favorite product"})
	}

}

func UnFav(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var json dto.FavoriteRequest
	var user model.User
	var fav model.Favorite
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.Conn.Find(&user, userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	db.Conn.Find(&fav, "product_id =? and user_id=?", json.ProductID, userId)
	if fav.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Not Found"})
		return
	} else {
		db.Conn.Delete(&fav).Where("product_id =? and user_id=?", json.ProductID, userId)
		c.JSON(http.StatusOK, gin.H{"message": "Un Favorites ?"})
		return
	}

	// var product model.Product
	// db.Conn.Find(&product, "id =?", json.ProductID)
	// if product.ID == 0 {
	// 	c.JSON(http.StatusOK, gin.H{"My product": "Not Found"})
	// 	return
	// }
	// db.Conn.Preload("Product").Delete(&fav).Where("user_id =? and product_id=?", uint(userId), json.ProductID)
	// c.JSON(http.StatusOK, gin.H{"mes": fav})
	// result := dto.FavoriteReponse{
	// 	ProductID: fav.ProductID,
	// }
	// c.JSON(http.StatusOK, gin.H{"fav": result, "message": "un favorite"})
}

func MyFav(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var favs []model.Favorite
	db.Conn.Model(&model.Favorite{}).Preload("Product").Find(&favs, "user_id=?", uint(userId))
	var result []dto.ProductRead
	for _, product := range favs {
		result = append(result, dto.ProductRead{
			ID:        product.ProductID,
			Name:      product.Product.Name,
			Desc:      product.Product.Description,
			Available: product.Product.Available,
			Price:     product.Product.Price,
			Weight:    product.Product.Weight,
			Image:     product.Product.Image,
			Rating:    product.Product.Rating,
		})
	}
	c.JSON(http.StatusOK, result)
}
