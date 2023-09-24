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

func DeleteStore(c *gin.Context) {
	//adminId := c.MustGet("storeId").(float64)
	var user model.Store
	var json Userid
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.Conn.Find(&user, "id =?", json.Id)
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Not Found"})
		return
	} else {
		db.Conn.Delete(&user).Where("id =?", json.Id)
		c.JSON(http.StatusOK, gin.H{"message": "Delete store ?"})
		return
	}
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
}

func UpdateStore(c *gin.Context) {
	id := c.Param("id")
	var store model.Store
	var json StoreUpdate
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&store, "id =?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
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

	db.Conn.Model(&store).Updates(StoreUpdate{UserName: json.UserName, Tel: json.Tel,
		NameStore: json.NameStore, Address: json.Address, District: json.District, SubDistrict: json.SubDistrict,
		Province: json.Province, Zipcode: json.Zipcode,
		AccountNumber: json.AccountNumber, AccountName: json.AccountName, Bank: json.Bank})

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "update store", "store": store})

}

func ReadAllStore(c *gin.Context) {
	var stores []model.Store
	db.Conn.Find(&stores)
	c.JSON(http.StatusOK, stores)
}

func ReadOneStore(c *gin.Context) {
	id := c.Param("id")
	var store model.Store
	if err := db.Conn.Find(&store, "id =?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if store.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"store": "Not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "store Read Success", "store": store})
}
