package untoken

import (
	"RmuttPlace/db"
	"RmuttPlace/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShipmentAll(c *gin.Context) {
	var shipment []model.Shipment
	db.Conn.Find(&shipment)
	c.JSON(http.StatusOK, shipment)
}
