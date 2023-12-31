package store

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

var hmacSampleSecretStore []byte

func ReadAll(c *gin.Context) {
	var stores []model.Store
	db.Conn.Find(&stores)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "storeId": stores})
}

func Profile(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var store model.Store
	db.Conn.Find(&store, storeId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "storeId": store})
}

type StoreBody struct {
	UserName    string `form:"username"`
	Tel         string `form:"tel"`
	NameStore   string `form:"namestore"`
	Address     string `form:"address"`
	District    string `form:"district"`
	SubDistrict string `form:"subdistrict"`
	Province    string `form:"province"`
	Zipcode     string `form:"zipcode"`
}

func UpdateMyStore(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var store model.Store
	var json StoreBody
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&store, "id =?", int(storeId)).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/stores/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(store.Image)
		store.Image = imagePath
	}
	file, err := c.FormFile("file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if file != nil {
		filePath := "./uploads/files/" + uuid.New().String()
		c.SaveUploadedFile(file, filePath)
		os.Remove(store.File)
		store.File = filePath
	}

	db.Conn.Save(&store)

	db.Conn.Model(&store).Updates(StoreBody{UserName: json.UserName, Tel: json.Tel,
		NameStore: json.NameStore, Address: json.Address, District: json.District, SubDistrict: json.SubDistrict,
		Province: json.Province, Zipcode: json.Zipcode})

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "update store", "store": store})

}

// func ResetPressStore(c *gin.Context) {}
