package user

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewBody struct {
	Comment string
	Rating  float32
}
type ReviewBodyRead struct {
	ProductID int
	UserID    int
	Comment   string
	Rating    float32
}

func CreateReview(c *gin.Context) {
	id := c.Param("id")
	userId := c.MustGet("userId").(float64)
	var json ReviewBody
	var product model.Product
	var review model.Review

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := db.Conn.Find(&product, "id = ?", id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query2 := db.Conn.Preload("User").Find(&review, "product_id=?", id)
	if err := query2.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var count int64
	var s int
	if err := db.Conn.Model(&model.Review{}).Select("sum(rating) as s").Where("product_id = ?", id).Count(&count).First(&s).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// n, _ := strconv.Atoi(id)
	// product.Rating = (float32(s) + json.Rating) / (float32(count) + 1)
	review.ProductID = int(product.ID)
	review.Comment = json.Comment
	review.Rating = json.Rating
	review.UserID = int(userId)
	db.Conn.Create(&review)
	// db.Conn.Save(&product)
	result := ReviewBodyRead{
		ProductID: review.ProductID,
		UserID:    review.UserID,
		Comment:   review.Comment,
		Rating:    review.Rating,
	}
	c.JSON(http.StatusOK, result)

	c.JSON(http.StatusOK, gin.H{"test": userId})
}
