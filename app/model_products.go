package app

import (
	"errors"
	"fmt"
	"log"
	"time"

	db "github.com/chuongthanh0410/interview/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// /////////////////
// modelA request
// /////////////////
type ModelProduct struct {
	b    Basemodel
	name string
}

type Product struct {
	// @format uuid
	ID        uuid.UUID `json:"id"`
	Reference string    `json:"reference"`
	Name      string    `json:"name"`
	// @Description Date when the product was added
	// @example "2025-03-22" "YYYY-MM-DD"
	// @format date
	AddedDate  time.Time `json:"added_date"`
	Status     string    `json:"status"`
	// @format uuid
	CategoryID uuid.UUID `json:"category_id"`
	Price      float64   `json:"price"`
	StockCity  string    `json:"stock_city"`
	// @format uuid
	SupplierID uuid.UUID `json:"supplier_id"`
	Quantity   int       `json:"quantity"`
}

func init() {
	m := &ModelProduct{
		name: "products",
	}
	RegisterModel(m, m.name)
}

func (m *ModelProduct) Authentication(c *gin.Context) (bool, error) {
	return m.b.Authentication(c)
}

func (m *ModelProduct) BeforeRequest(c *gin.Context) (bool, error) {
	return m.b.BeforeRequest(c)
}

func (m *ModelProduct) Process(c *gin.Context, function_id string) (bool, error) {
	fmt.Println("[ModelProduct-model] process ...", c.ClientIP())
	switch function_id {
	case "get_id":
		return m.GetOne(c)

	case "update_id":
		return m.UpdateOne(c)

	case "get_all":
		return m.GetMany(c)

	case "update_many":
		return m.UpdateMany(c)

	case "delete_id":
		return m.DeleteOne(c)

	default:
		c.JSON(200, gin.H{"message": "Unknown handler!!"})
	}
	return true, nil
}

func (m *ModelProduct) AfterRequest(c *gin.Context) (bool, error) {
	fmt.Println("[ModelProduct-model] afterRequest auditlog ...", c.ClientIP())
	return true, nil
}

// @host localhost:8080
// @BasePath /api

// @Tags products
// @Summary Get one products
// @Description Get one products
// @Produce json
// @Success 200 {array} Product
// @Router /products/:id [get]
func (m *ModelProduct) GetOne(c *gin.Context) (bool, error) {
	id := c.Param("id")
	rows, err := db.DBClient.SelectOne(m.name, "id", id)

	if err != nil {
		log.Printf(err.Error())
		return m.b.ReturnErr(c, fmt.Sprintf("get product %s error", id))
	}

	defer rows.Close()
	products, err := m.ScanData(m.name, []string{}, rows)
	if err != nil {
		return m.b.ReturnErr(c, fmt.Sprintf("get product %s error %s", id, err.Error()))

	}

	return m.b.ReturnSuc(c, products, nil)
}

func (m *ModelProduct) ScanData(table string, getFields []string, rows pgx.Rows) (interface{}, error) {
	var products []Product
	for rows.Next() {
		product := Product{}
		err := rows.Scan(
			&product.ID,
			&product.Reference,
			&product.Name,
			&product.AddedDate,
			&product.Status,
			&product.CategoryID,
			&product.Price,
			&product.StockCity,
			&product.SupplierID,
			&product.Quantity,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during iteration: %v", err)
		return nil, errors.New("Error iterating through rows")
	}
	return products, nil
}

func (m *ModelProduct) GetMany(c *gin.Context) (bool, error) {
	c.JSON(200, gin.H{"message": "using custom_search instead"})
	return true, nil
}

func (m *ModelProduct) UpdateOne(c *gin.Context) (bool, error) {
	return m.b.UpdateOne(c, m.name)
}

func (m *ModelProduct) UpdateMany(c *gin.Context) (bool, error) {
	// TODO
	c.JSON(400, gin.H{"message": "Not implement yet"})
	return true, nil
}

func (m *ModelProduct) DeleteOne(c *gin.Context) (bool, error) {
	return m.b.DeleteOne(c, m.name)
}
