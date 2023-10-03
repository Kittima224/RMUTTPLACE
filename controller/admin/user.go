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

func ReadAllUser(c *gin.Context) {
	var users []model.User
	db.Conn.Find(&users)
	c.JSON(http.StatusOK, users)
}

func ReadOneUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := db.Conn.Find(&user, "id =?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"user": "Not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user})
}

type Userid struct {
	Id uint `json:"id"`
}

func DeleteUser(c *gin.Context) {
	var user model.User
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
		c.JSON(http.StatusOK, gin.H{"message": "Delete user ?"})
		return
	}
}

type UpdateUserByAdmin struct {
	UserName    string `form:"username"`
	Tel         string `form:"tel"`
	Address     string `form:"address"`
	District    string `form:"district"`
	SubDistrict string `form:"subdistrict"`
	Province    string `form:"province"`
	Zipcode     string `form:"zipcode"`
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	var json UpdateUserByAdmin
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Conn.Find(&user, "id =?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/profileusers/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(user.Image)
		user.Image = imagePath
	}
	db.Conn.Save(&user)
	db.Conn.Model(&user).Updates(UpdateUserByAdmin{UserName: json.UserName, Tel: json.Tel,
		Address: json.Address, District: json.District, SubDistrict: json.SubDistrict,
		Province: json.Province, Zipcode: json.Zipcode})

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "update user", "user": user})
}

func UpdateUserPhoto(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := db.Conn.Find(&user, "id =?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/profileusers/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(user.Image)
		user.Image = imagePath
	}
	db.Conn.Save(&user)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "update user", "user": user})
}

type CreateUserByAdminRequest struct {
	Email       string `form:"email" binding:"required"`
	Password    string `form:"password" binding:"required"`
	UserName    string `form:"username"`
	Tel         string `form:"tel"`
	Address     string `form:"address"`
	District    string `form:"district"`
	SubDistrict string `form:"subdistrict"`
	Province    string `form:"province"`
	Zipcode     string `form:"zipcode"`
}

func Register(c *gin.Context) {
	var json CreateUserByAdminRequest
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check user exists
	var userExist model.User
	db.Conn.Where("email = ?", json.Email).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "User Exists",
		})
		return
	}
	encrytedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := model.User{Email: json.Email, UserName: json.UserName, Password: string(encrytedPassword), Tel: json.Tel,
		Address: json.Address, District: json.District, SubDistrict: json.SubDistrict, Province: json.Province,
		Zipcode: json.Zipcode}
	db.Conn.Create(&user)
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/profileusers/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(user.Image)
		user.Image = imagePath
	}
	db.Conn.Save(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "user create success",
			"userId":  user.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "user create failed",
		})
	}
}
