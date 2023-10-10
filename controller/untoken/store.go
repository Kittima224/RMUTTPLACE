package untoken

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
