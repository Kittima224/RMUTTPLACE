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

var hmacSampleSecretAdmin []byte

func Profile(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	db.Conn.Find(&admin, adminId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": admin})
}

func ReadAll(c *gin.Context) {
	var admins []model.Admin
	db.Conn.Find(&admins)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "admin": admins})
}
func GetProfile(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	db.Conn.Find(&admin, adminId)
	c.JSON(http.StatusOK, gin.H{"adminImage": admin.Image})
}

type AdminBody struct {
	UserName string `form:"username"`
}

func UpdateAdmin(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	var json AdminBody
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
		os.Remove(admin.Image)
		admin.Image = imagePath
	}
	admin.UserName = json.UserName
	db.Conn.Save(&admin)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "update product", "product": admin})
}
