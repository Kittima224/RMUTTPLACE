package untoken

import (
	"RmuttPlace/db"
	"RmuttPlace/dto"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductAllStore(c *gin.Context) {
	id := c.Param("id")
	var products []model.Product
	query := db.Conn.Preload("Category").Find(&products, "store_id", id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var result []dto.ProductRead
	for _, product := range products {
		result = append(result, dto.ProductRead{
			ID:        product.ID,
			Name:      product.Name,
			Desc:      product.Desc,
			Available: product.Available,
			Price:     product.Price,
			Weight:    product.Weight,
			Image:     product.Image,
			Category: model.CategoryRead{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
		})
	}
	c.JSON(http.StatusOK, result)

}
func ReadProductAll(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")
	var products []model.Product
	db.Conn.Preload("Category").Find(&products)

	if category != "" {
		db.Conn.Find(&products, "category LIKE ?", "%"+category+"%")
	}
	if search != "" {
		db.Conn.Find(&products, "name LIKE ? or desc LIKE? ", "%"+search+"%")
	}
	var result []dto.ProductRead
	for _, product := range products {
		result = append(result, dto.ProductRead{
			ID:        product.ID,
			Name:      product.Name,
			Desc:      product.Desc,
			Available: product.Available,
			Price:     product.Price,
			Weight:    product.Weight,
			Image:     product.Image,
			Category: model.CategoryRead{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
		})
	}
	c.JSON(http.StatusOK, result)
}
