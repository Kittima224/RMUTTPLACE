package admin

import (
	"RmuttPlace/db"
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
		imagePath := "./uploads/pd/" + uuid.New().String()
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

type ProductRead struct {
	ID        uint
	Name      string
	Desc      string
	Category  model.CategoryRead
	Available int
	Price     int
	Weight    int
}

func ReadOneProduct(c *gin.Context) {
	id := c.Param("id")
	var product model.Product
	query := db.Conn.Preload("Images").Preload("Category").Find(&product, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "product Read Success", "product": product})
}

func ReadProductAllMyStore(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var products []model.Product
	db.Conn.Find(&products, "store_id", storeId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "products": products})
}
func FindNameProduct(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")
	var products []model.Product
	if category != "" {
		db.Conn.Find(&products, "category LIKE ?", "%"+category+"%")
	}
	if search != "" {
		db.Conn.Find(&products, "name LIKE ? ", "%"+search+"%")
	}
	c.JSON(http.StatusOK, gin.H{"products": products})

	//ใช้เก็บข้อมูลมาวิเคราะห์
}
func ReadProductAll(c *gin.Context) {
	var products []model.Product
	db.Conn.Preload("Category").Find(&products)

	var result []ProductRead
	for _, product := range products {
		result = append(result, ProductRead{
			ID:        product.ID,
			Name:      product.Name,
			Desc:      product.Desc,
			Available: product.Available,
			Price:     product.Price,
			Weight:    product.Weight,
			Category: model.CategoryRead{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
		})
	}
	c.JSON(http.StatusOK, result)
}
