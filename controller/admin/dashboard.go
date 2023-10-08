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

	// type PieChart struct {
	// 	Name  string
	// 	Value int
	// }
	// data := []PieChart
	// rows, err := db.Conn.Query("Select name,count(name)  as value, from categories group by name")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer rows.Close()
	// var name string

	// for rows.Next() {
	// 	err := rows.Scan(&PieChart{Name: name})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	data = append(data, model.Category{Name: name})
	// }
	// fmt.Println(data)
	// type pieRead struct {
	// 	ID    uint   `json:"id"`
	// 	Name  string `json:"name"`
	// 	Value int    `json:"value"`
	// }
	// var result []pieRead

	// var id uint
	// var name string
	// var value int
	//db.Conn.Raw("SELECT categories.id as id,categories.name as name ,COUNT(products.id) as value from products JOIN categories on products.category_id = categories.id WHERE products.deleted_at is null GROUP by products.category_id").Scan(&result)
	// result = append(result, pieRead{
	// 	ID:    id,
	// 	Name:  name,
	// 	Value: value,
	// })

	// type chart struct {
	// 	ID    uint   `json:"id"`
	// 	Name  string `json:"name"`
	// 	Value int    `json:"value"`
	// }
	// var r []chart
	// ro := db.Conn.Raw("SELECT DATE_FORMAT(ot.created_at,'%M') as name,SUM(ot.quantity*p.price) as value,DATE_FORMAT(ot.created_at,'%m') as id FROM order_items as ot JOIN products as p ON ot.product_id=p.id GROUP BY DATE_FORMAT(created_at,'%M')").Row()
	// ro.Scan(&name, &value, &id)
	// r = append(r, chart{
	// 	ID:    id,
	// 	Name:  name,
	// 	Value: value,
	// })
	c.JSON(http.StatusOK, gin.H{"total_price": humanize.Commaf(float64(totalPrice)),
		"count_store":   humanize.Commaf(float64(cstore)),
		"count_acc":     humanize.Commaf(float64(cuser)),
		"count_product": humanize.Commaf(float64(cproduct)),
		// "pie":           re,
		// "chart":         r,
	})
}

type ChartResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}
type PieChart struct {
	Name  string
	Value int
}
