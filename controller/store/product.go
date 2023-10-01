package store

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductBody struct {
	Name       string `form:"name" binding:"required"`
	Desc       string `form:"desc"`
	CategoryID int    `form:"categoryId" binding:"required"`
	Available  int    `form:"available" binding:"required"`
	Price      int    `form:"price" binding:"required"`
	Weight     int    `form:"weight" binding:"required"`
}

func Create(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var store model.Store
	query := db.Conn.Find(&store, storeId)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var json ProductBody
	if err := c.ShouldBindWith(&json, binding.FormMultipart); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	product := model.Product{Name: json.Name, Desc: json.Desc, StoreID: int(storeId),
		CategoryID: json.CategoryID,
		Available:  json.Available, Price: json.Price,
		Weight: json.Weight}
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/pd/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		product.Image = imagePath
	}

	db.Conn.Create(&product)

	if product.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"message":   "product create success",
			"productId": product.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "product create failed",
		})
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

type ProductUpdateBody struct {
	Name       string `form:"name"`
	Desc       string `form:"desc"`
	CategoryID string `form:"categoryId"`
	Available  int    `form:"available"`
	Price      int    `form:"price"`
	Weight     int    `form:"weight"`
}

func UpdateProductMystore(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var product model.Product
	var json ProductUpdateBody
	if err := c.ShouldBindWith(&json, binding.FormMultipart); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := db.Conn.Find(&product, "store_id = ? AND id = ?", int(storeId), id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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

	c.JSON(http.StatusOK, gin.H{"My product": storeId, "productid": product})
}

func DeleteProductMyStore(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var product model.Product
	db.Conn.Find(&product, "store_id =? and id =?", int(storeId), id)
	if product.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"My product": "Not Found"})
		return
	} else {
		db.Conn.Delete(&product).Where("store_id =? and id =?", int(storeId), id)
		c.JSON(http.StatusOK, gin.H{"My product": "Delete success"})
		return
	}
}

func FindOneProductMyStore(c *gin.Context) {
	id := c.Param("id")
	storeId := c.MustGet("storeId").(float64)
	var product model.Product
	var store model.Store
	query := db.Conn.Preload("Images").Find(&product, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&store, "id =?", storeId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})

}

func ReadProductAllMyStore(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var store model.Store
	var products []model.Product
	db.Conn.Find(&products, "store_id =?", int(storeId))

	if err := db.Conn.Find(&store, "id =?", int(storeId)).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
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
	db.Conn.Find(&products)
	c.JSON(http.StatusOK, products)
}
