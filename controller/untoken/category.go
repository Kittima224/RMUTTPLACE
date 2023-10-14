package untoken

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryAll(c *gin.Context) {
	var categories []model.Category
	if err := db.Conn.Find(&categories).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var result []model.CategoryRead
	for _, r := range categories {
		result = append(result, model.CategoryRead{
			ID:   r.ID,
			Name: r.Name,
		})
	}
	c.JSON(http.StatusOK, result)
}
