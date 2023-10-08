package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/dto"
	"RmuttPlace/model"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductUpdateBody struct {
	Name       string `form:"name"`
	Desc       string `form:"desc"`
	CategoryID uint   `form:"categoryId"`
	Available  int    `form:"available"`
	Price      int    `form:"price"`
	Weight     int    `form:"weight"`
}

func UpdateProduct(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	var product model.Product
	var json ProductUpdateBody
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&product, "id =?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/products/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(product.Image)
		product.Image = imagePath
	}
	db.Conn.Save(&product)
	db.Conn.Model(&product).Updates(ProductUpdateBody{Name: json.Name,
		Desc: json.Desc, CategoryID: json.CategoryID, Available: json.Available, Price: json.Price, Weight: json.Weight})

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "update product", "product": product})
}

func DeleteProduct(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var product model.Product
	var json Userid
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.Conn.Find(&product, "id =?", json.Id)
	if product.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Not Found"})
		return
	} else {
		db.Conn.Delete(&product).Where("id =?", json.Id)
		c.JSON(http.StatusOK, gin.H{"message": "Delete product ?"})
		return
	}
}

func ReadOneProduct(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
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

func ReadProductAllMyStore(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var products []model.Product
	db.Conn.Find(&products, "store_id", storeId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "products": products})
}

func ReadProductAll(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var products []model.Product
	db.Conn.Preload("Category").Find(&products)

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
