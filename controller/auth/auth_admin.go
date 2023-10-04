package auth

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecretAdmin []byte

type RegisterAdmineBody struct {
	Email    string `form:"email" binding:"required"`
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func RegisterAdmin(c *gin.Context) {
	var json RegisterAdmineBody
	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check Admin exists
	var adminExist model.Admin
	db.Conn.Where("email = ?", json.Email).First(&adminExist)
	if adminExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Addmin Exists",
		})
		return
	}
	//create Addmin
	encrytedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	admin := model.Admin{Email: json.Email, UserName: json.UserName, Password: string(encrytedPassword)}
	db.Conn.Create(&admin)
	image, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if image != nil {
		imagePath := "./uploads/admins/" + uuid.New().String()
		c.SaveUploadedFile(image, imagePath)
		os.Remove(admin.Image)
		admin.Image = imagePath
	}
	db.Conn.Save(&admin)
	if admin.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "admin create success",
			"adminId": admin.ID,
			"image":   adminExist.Image,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "admin create failed",
		})
	}
}

type LoginAdminBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginAdmin(c *gin.Context) {
	var json LoginAdminBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check admin exists
	var adminExist model.Admin
	db.Conn.Where("email = ?", json.Email).First(&adminExist)
	if adminExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Admin Does Not Exists",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(adminExist.Password), []byte(json.Password))
	if err == nil {
		hmacSampleSecretAdmin = []byte(os.Getenv("JWT_SECRET_KEY_ADMIN"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"adminId": adminExist.ID,
			"exp":     time.Now().Add(time.Hour * 1).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecretAdmin)
		fmt.Println(tokenString, err)
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Login success",
			"token":   tokenString,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Loging Failed",
		})
	}
}
