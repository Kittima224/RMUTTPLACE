package store

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	UserName      string `form:"username" binding:"required"`
	Tel           string `form:"tel" binding:"required"`
	NameStore     string `form:"namestore" binding:"required"`
	Address       string `form:"address" binding:"required"`
	District      string `form:"district" binding:"required"`
	SubDistrict   string `form:"subdistrict" binding:"required"`
	Province      string `form:"province" binding:"required"`
	Zipcode       string `form:"zipcode" binding:"required"`
	Shipment      string `form:"shipmentname[]" binding:"required"`
	AccountNumber string `form:"account_number" binding:"required"`
	AccountName   string `form:"account_name" binding:"required"`
	Bank          string `form:"bank" binding:"required"`
}

func AddProfileMystore(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var store model.Store
	var json StoreBody
	if err := c.ShouldBindWith(&json, binding.FormMultipart); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := db.Conn.Find(&store, "id =?", int(storeId)).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// form, _ := c.MultipartForm()
	// files := form.Value["shipmentname[]"]
	// var ships []model.Shipment
	// for _, file := range files {
	// 	ships = append(ships, model.Shipment{
	// 		StoreID:      uint(storeId),
	// 		ShipmentName: file,
	// 	})
	// }
	// store.Shipments = ships

	db.Conn.Updates(&store)

	db.Conn.Model(&store).Updates(StoreUpdate{UserName: json.UserName, Tel: json.Tel,
		NameStore: json.NameStore, Address: json.Address, District: json.District, SubDistrict: json.SubDistrict,
		Province: json.Province, Zipcode: json.Zipcode,
		AccountNumber: json.AccountNumber, AccountName: json.AccountName, Bank: json.Bank})

	c.JSON(http.StatusOK, gin.H{"My store": store})
}

type StoreUpdate struct {
	UserName      string `form:"username"`
	Tel           string `form:"tel"`
	NameStore     string `form:"namestore"`
	Address       string `form:"address"`
	District      string `form:"district"`
	SubDistrict   string `form:"subdistrict"`
	Province      string `form:"province"`
	Zipcode       string `form:"zipcode"`
	AccountNumber string `form:"account_number"`
	AccountName   string `form:"account_name"`
	Bank          string `form:"bank"`
	//image         *multipart.FileHeader `form:"image"`
}

func UpdateProfileMystore(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var store model.Store
	var json StoreUpdate
	if err := c.ShouldBindWith(&json, binding.FormMultipart); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := db.Conn.Find(&store, "id =?", int(storeId)).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/profilestores/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(store.Image)
		store.Image = imagePath
	}
	db.Conn.Save(&store)

	db.Conn.Save(&store)

	db.Conn.Model(&store).Updates(StoreUpdate{UserName: json.UserName, Tel: json.Tel,
		NameStore: json.NameStore, Address: json.Address, District: json.District, SubDistrict: json.SubDistrict,
		Province: json.Province, Zipcode: json.Zipcode,
		AccountNumber: json.AccountNumber, AccountName: json.AccountName, Bank: json.Bank})

	c.JSON(http.StatusOK, gin.H{"My store": store})
}

type ChangeProfileStoretBody struct {
	Image *multipart.FileHeader `form:"image"`
}

func ChangProfileMystore(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	var store model.Store
	var json ChangeProfileStoretBody
	if err := c.ShouldBindWith(&json, binding.FormMultipart); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := db.Conn.Find(&store, "id =?", uint(storeId)).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/profilestores/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(store.Image)
		store.Image = imagePath
	}
	db.Conn.Save(&store)
	c.JSON(http.StatusOK, gin.H{"My store": store})
}

func ResetPressStore(c *gin.Context) {}
