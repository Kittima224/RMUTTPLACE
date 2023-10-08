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
			Rating: product.Rating,
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

// แก้
func FindOneProductMyStore(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var product model.Product
	var store model.Store
	var reviews []model.Review
	query := db.Conn.Preload("Store").Preload("Category").Find(&product, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	query2 := db.Conn.Preload("User").Find(&reviews, "product_id", id)
	if err := query2.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&store, "id =?", storeId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	result := dto.ProductReadOne{
		ID:        product.ID,
		Name:      product.Name,
		Desc:      product.Desc,
		Available: product.Available,
		Image:     product.Image,
		Price:     product.Price,
		Weight:    product.Weight,
		Category: model.CategoryRead{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		},
		Store: dto.StoreRead{
			ID:   product.Store.ID,
			Name: product.Store.NameStore,
		},
	}
	var rv []dto.ReviewBodyRead
	for _, r := range reviews {
		rv = append(rv, dto.ReviewBodyRead{
			UserID:  r.UserID,
			Name:    r.User.UserName,
			Comment: r.Comment,
			Rating:  r.Rating,
		})
	}
	result.Reviews = rv
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "product Read Success", "product": result})

}
