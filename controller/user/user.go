package user

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

var hmacSampleSecret []byte

func ReadAll(c *gin.Context) {
	var users []model.User
	db.Conn.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "userId": users})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user model.User
	db.Conn.Find(&user, userId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "userId": user})
}

type UserBody struct {
	UserName    string                `form:"username"  binding:"required"`
	Tel         string                `form:"tel"  binding:"required"`
	Address     string                `form:"address" binding:"required"`
	District    string                `form:"district" binding:"required"`
	SubDistrict string                `form:"subdistrict" binding:"required"`
	Province    string                `form:"province" binding:"required"`
	Zipcode     string                `form:"zipcode" binding:"required"`
	image       *multipart.FileHeader `form:"image"`
}

func AddProfileUser(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user model.User
	var json UserBody
	if err := c.ShouldBindWith(&json, binding.FormMultipart); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := db.Conn.Find(&user, "id =?", uint(userId)).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, _ := c.FormFile("image")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	if image != nil {
		imagePath := "./uploads/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(user.Image)
		user.Image = imagePath
	}
	db.Conn.Save(&user)
	db.Conn.Model(&user).Updates(UserBody{UserName: json.UserName, Tel: json.Tel,
		Address: json.Address, District: json.District, SubDistrict: json.SubDistrict,
		Province: json.Province, Zipcode: json.Zipcode})

	c.JSON(http.StatusOK, gin.H{"My user": user})
}

type UpdateUser struct {
	UserName    string                `form:"username"`
	Tel         string                `form:"tel"`
	Address     string                `form:"address"`
	District    string                `form:"district"`
	SubDistrict string                `form:"subdistrict"`
	Province    string                `form:"province"`
	Zipcode     string                `form:"zipcode"`
	image       *multipart.FileHeader `form:"image"`
}

func UpdateProfileUser(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user model.User
	var json UpdateUser
	if err := c.ShouldBindWith(&json, binding.FormMultipart); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := db.Conn.Find(&user, "id =?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(user.Image)
		user.Image = imagePath
	}
	db.Conn.Save(&user)
	db.Conn.Model(&user).Updates(UpdateUser{UserName: json.UserName, Tel: json.Tel,
		Address: json.Address, District: json.District, SubDistrict: json.SubDistrict,
		Province: json.Province, Zipcode: json.Zipcode})

	c.JSON(http.StatusOK, gin.H{"My user": user})
}

func ResetPasswordUser(c *gin.Context) {}
func DeleteUser(c *gin.Context)        {}
