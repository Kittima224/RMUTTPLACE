package auth

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecretStore []byte

type RegisterStoreBody struct {
	Email    string `json:"email" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterStore(c *gin.Context) {
	var json RegisterStoreBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check Store exists
	var storeExist model.Store
	db.Conn.Where("email = ?", json.Email).First(&storeExist)
	if storeExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Store Does Not Exists",
		})
		return
	}
	//create Store
	encrytedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	store := model.Store{Email: json.Email, UserName: json.UserName, Password: string(encrytedPassword)}
	db.Conn.Create(&store)
	if store.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "store create success",
			"userId":  store.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "store create failed",
		})
	}
}

type LoginStoreBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginStore(c *gin.Context) {
	var json LoginStoreBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check user exists
	var storeExist model.Store
	db.Conn.Where("email = ?", json.Email).First(&storeExist)
	if storeExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Store Does Not Exists",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(storeExist.Password), []byte(json.Password))
	if err == nil {
		hmacSampleSecretStore = []byte(os.Getenv("JWT_SECRET_KEY_STORE"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"storeId": storeExist.ID,
			"exp":     time.Now().Add(time.Hour * 1).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecretStore)
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
