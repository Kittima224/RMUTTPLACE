package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func DeleteStore(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
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
	UserName    string `form:"username"`
	Tel         string `form:"tel"`
	NameStore   string `form:"namestore"`
	Address     string `form:"address"`
	District    string `form:"district"`
	SubDistrict string `form:"subdistrict"`
	Province    string `form:"province"`
	Zipcode     string `form:"zipcode"`
	Status      bool   `form:"status"`
}

func UpdateStore(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
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
		imagePath := "./uploads/stores/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(store.Image)
		store.Image = imagePath
	}

	db.Conn.Save(&store)

	db.Conn.Model(&store).Updates(StoreUpdate{UserName: json.UserName, Tel: json.Tel,
		NameStore: json.NameStore, Address: json.Address, District: json.District, SubDistrict: json.SubDistrict,
		Province: json.Province, Zipcode: json.Zipcode, Status: json.Status,
	})

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "update store", "store": store})

}

func ReadAllStore(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var stores []model.Store
	db.Conn.Find(&stores)
	c.JSON(http.StatusOK, stores)
}

func ReadOneStore(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
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

type CreateStoreByAdminRequest struct {
	Email       string `form:"email" binding:"required"`
	Password    string `form:"password" binding:"required"`
	UserName    string `form:"username"`
	NameStore   string `form:"namestore"`
	Tel         string `form:"tel"`
	Address     string `form:"address"`
	District    string `form:"district"`
	SubDistrict string `form:"subdistrict"`
	Province    string `form:"province"`
	Zipcode     string `form:"zipcode"`
}

func StoreRegister(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var json CreateStoreByAdminRequest
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check user exists
	var storeExist model.Store
	db.Conn.Where("email = ?", json.Email).First(&storeExist)
	if storeExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Store Exists",
		})
		return
	}
	encrytedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	store := model.Store{Email: json.Email, UserName: json.UserName, Password: string(encrytedPassword), Tel: json.Tel,
		Address: json.Address, District: json.District, SubDistrict: json.SubDistrict, Province: json.Province,
		Zipcode: json.Zipcode, NameStore: json.NameStore}
	db.Conn.Create(&store)
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
	db.Conn.Save(&store)
	if store.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "store create success",
			"storeId": store.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "store create failed",
		})
	}
}
