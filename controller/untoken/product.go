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
			Desc:      product.Description,
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
	categoryid := c.Query("categoryid")
	desc := c.Query("desc")
	var products []model.Product
	var p []int
	db.Conn.Raw("SELECT products.id as p FROM products JOIN stores on products.store_id=stores.id WHERE stores.status=true").Scan(&p)
	query := db.Conn.Preload("Category").Preload("Store")
	if categoryid != "" {
		query = query.Where("category_id=?", categoryid)
	}
	if search != "" {

		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if desc != "" {
		query = query.Where("description like ?", "%"+desc+"%")
	}
	query.Find(&products, p)
	var result []dto.ProductRead
	for _, product := range products {
		result = append(result, dto.ProductRead{
			ID:        product.ID,
			Name:      product.Name,
			Desc:      product.Description,
			Available: product.Available,
			Price:     product.Price,
			Weight:    product.Weight,
			Image:     product.Image,
			Category: model.CategoryRead{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
			Rating: product.Rating,
			Store: dto.StoreRead{
				ID:   product.Store.ID,
				Name: product.Store.NameStore,
			},
		})
	}

	c.JSON(http.StatusOK, result)
}

func FindOneProduct(c *gin.Context) {
	id := c.Param("id")
	var product model.Product
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

	result := dto.ProductReadOne{
		ID:        product.ID,
		Name:      product.Name,
		Desc:      product.Description,
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
		Rating: product.Rating,
	}
	var rv []dto.ReviewBodyRead
	for _, r := range reviews {
		rv = append(rv, dto.ReviewBodyRead{
			UserID:  r.UserID,
			Comment: r.Comment,
			Rating:  r.Rating,
			User: dto.UserReview{
				ID:    r.User.ID,
				Name:  r.User.UserName,
				Image: r.User.Image,
			},
		})
	}
	result.Reviews = rv
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "product Read Success", "product": result})

}

func FifteenProduct(c *gin.Context) {
	var p []int
	var products []model.Product
	db.Conn.Raw("SELECT random(id) as p FROM products LIMIT 15").Scan(p)
	var result []dto.ProductRead
	db.Conn.Find(&products, p)
	for _, product := range products {
		result = append(result, dto.ProductRead{
			ID:        product.ID,
			Name:      product.Name,
			Desc:      product.Description,
			Available: product.Available,
			Price:     product.Price,
			Weight:    product.Weight,
			Image:     product.Image,
			Category: model.CategoryRead{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
			Rating: product.Rating,
			Store: dto.StoreRead{
				ID:   product.Store.ID,
				Name: product.Store.NameStore,
			},
		})
	}

	c.JSON(http.StatusOK, result)
}
func BestSeller(c *gin.Context) {
	var n []int
	var products []model.Product
	db.Conn.Raw("SELECT ot.product_id as n FROM order_items as ot  GROUP BY ot.product_id ORDER by COUNT(ot.product_id) DESC LIMIT 10").Scan(&n)
	var result []dto.ProductRead
	db.Conn.Preload("Store").Find(&products, n)
	for _, product := range products {
		result = append(result, dto.ProductRead{
			ID:        product.ID,
			Name:      product.Name,
			Desc:      product.Description,
			Available: product.Available,
			Price:     product.Price,
			Weight:    product.Weight,
			Image:     product.Image,
			Category: model.CategoryRead{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
			Rating: product.Rating,
			Store: dto.StoreRead{
				ID:   product.Store.ID,
				Name: product.Store.NameStore,
			},
		})
	}

	c.JSON(http.StatusOK, result)

}

func NotSeller(c *gin.Context) {
	var n []int
	var products []model.Product
	db.Conn.Raw("SELECT ot.product_id as n FROM order_items as ot  GROUP BY ot.product_id ORDER by COUNT(ot.product_id) LIMIT 10").Scan(&n)
	var result []dto.ProductRead
	db.Conn.Preload("Store").Find(&products, n)
	for _, product := range products {
		result = append(result, dto.ProductRead{
			ID:        product.ID,
			Name:      product.Name,
			Desc:      product.Description,
			Available: product.Available,
			Price:     product.Price,
			Weight:    product.Weight,
			Image:     product.Image,
			Category: model.CategoryRead{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
			Rating: product.Rating,
			Store: dto.StoreRead{
				ID:   product.Store.ID,
				Name: product.Store.NameStore,
			},
		})
	}

	c.JSON(http.StatusOK, result)
}
