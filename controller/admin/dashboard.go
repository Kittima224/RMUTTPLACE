package admin

import (

	// "fmt"
	"RmuttPlace/db"
	"RmuttPlace/model"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {
	adminId := c.MustGet("adminId").(float64)
	var admin model.Admin
	if err := db.Conn.Find(&admin, adminId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var totalPrice int
	db.Conn.Raw("SELECT sum(order_items.quantity*products.price) from order_items JOIN products ON products.id = order_items.product_id").Scan(&totalPrice)

	var cstore int
	db.Conn.Raw("SELECT COUNT(id) from stores where deleted_at is null").Scan(&cstore)

	var cuser int
	db.Conn.Raw("SELECT COUNT(id) from users where deleted_at is null").Scan(&cuser)

	var cproduct int
	db.Conn.Raw("SELECT COUNT(id) from products WHERE deleted_at is null").Scan(&cproduct)

	type Chart struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	var r []Chart
	db.Conn.Raw("SELECT DATE_FORMAT(ot.created_at,'%b') as name,SUM(ot.quantity*p.price) as value,DATE_FORMAT(ot.created_at,'%c') as id FROM order_items as ot JOIN products as p ON ot.product_id=p.id GROUP BY DATE_FORMAT(created_at,'%b');").Scan(&r)

	var p []Chart
	db.Conn.Raw("SELECT categories.id as id,categories.name as name ,COUNT(products.id) as value from products JOIN categories on products.category_id = categories.id WHERE products.deleted_at is null GROUP by products.category_id").Scan(&p)

	c.JSON(http.StatusOK, gin.H{
		"total_price":   humanize.Commaf(float64(totalPrice)),
		"count_store":   humanize.Commaf(float64(cstore)),
		"count_acc":     humanize.Commaf(float64(cuser)),
		"count_product": humanize.Commaf(float64(cproduct)),
		"pie":           p,
		"chart":         r,
	})
}

type YearOfChart struct {
	Year string `json:"year"`
}

func DashboardTest(c *gin.Context) {
	var totalPrice int
	db.Conn.Raw("SELECT sum(order_items.quantity*products.price) from order_items JOIN products ON products.id = order_items.product_id").Scan(&totalPrice)

	var cstore int
	db.Conn.Raw("SELECT COUNT(id) from stores where deleted_at is null").Scan(&cstore)

	var cuser int
	db.Conn.Raw("SELECT COUNT(id) from users where deleted_at is null").Scan(&cuser)

	var cproduct int
	db.Conn.Raw("SELECT COUNT(id) from products WHERE deleted_at is null").Scan(&cproduct)

	// var json YearOfChart
	// if err := c.ShouldBindJSON(&json); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	type Chart struct {
		ID int    `json:"id"`
		X  string `json:"x"`
		Y  int    `json:"y"`
	}
	var r []Chart
	db.Conn.Raw("SELECT DATE_FORMAT(ot.created_at,'%b') as x,SUM(ot.quantity*p.price) as y,DATE_FORMAT(ot.created_at,'%c') as id FROM order_items as ot JOIN products as p on ot.product_id=p.id GROUP BY x ORDER BY ot.id").Scan(&r)

	type Pie struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	var p []Pie
	db.Conn.Raw("SELECT categories.id as id,categories.name as name ,COUNT(products.id) as value from products JOIN categories on products.category_id = categories.id WHERE products.deleted_at is null GROUP by products.category_id").Scan(&p)

	c.JSON(http.StatusOK, gin.H{
		"total_price":   humanize.Commaf(float64(totalPrice)),
		"count_store":   humanize.Commaf(float64(cstore)),
		"count_acc":     humanize.Commaf(float64(cuser)),
		"count_product": humanize.Commaf(float64(cproduct)),
		"pie":           p,
		"chart":         r,
	})
}
