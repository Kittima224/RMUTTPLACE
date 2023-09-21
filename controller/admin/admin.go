package admin

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"net/http"

	"github.com/gin-gonic/gin"
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
