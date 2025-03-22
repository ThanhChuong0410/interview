package app

import (
	"errors"
	"fmt"
	"log"

	db "github.com/chuongthanh0410/interview/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// /////////////////
// modelA request
// /////////////////
type ModelSuppliers struct {
	b    Basemodel
	name string
}

type Suppliers struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func init() {
	m := &ModelSuppliers{
		name: "suppliers",
	}
	RegisterModel(m, m.name)
}

func (m *ModelSuppliers) Authentication(c *gin.Context) (bool, error) {
	return m.b.Authentication(c)
}

func (m *ModelSuppliers) BeforeRequest(c *gin.Context) (bool, error) {
	return m.b.BeforeRequest(c)
}

func (m *ModelSuppliers) Process(c *gin.Context, function_id string) (bool, error) {
	fmt.Println("[ModelSuppliers-model] process ...", c.ClientIP())
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

func (m *ModelSuppliers) AfterRequest(c *gin.Context) (bool, error) {
	fmt.Println("[ModelSuppliers-model] afterRequest auditlog ...", c.ClientIP())
	return true, nil
}

func (m *ModelSuppliers) GetOne(c *gin.Context) (bool, error) {
	id := c.Param("id")
	rows, err := db.DBClient.SelectOne(m.name, "id", id)

	if err != nil {
		log.Printf(err.Error())
		return m.b.ReturnErr(c, fmt.Sprintf("get suppliers %s error", id))
	}

	defer rows.Close()
	suppliers, err := m.ScanData(m.name, []string{}, rows)
	if err != nil {
		return m.b.ReturnErr(c, fmt.Sprintf("get suppliers %s error %s", id, err.Error()))

	}

	return m.b.ReturnSuc(c, suppliers, nil)
}

func (m *ModelSuppliers) ScanData(table string, getFields []string, rows pgx.Rows) (interface{}, error) {
	var suppliers []Suppliers
	for rows.Next() {
		supplier := Suppliers{}
		err := rows.Scan(
			&supplier.ID,
			&supplier.Name,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during iteration: %v", err)
		return nil, errors.New("Error iterating through rows")
	}
	return suppliers, nil
}

func (m *ModelSuppliers) GetMany(c *gin.Context) (bool, error) {
	c.JSON(200, gin.H{"message": "using custom_search instead"})
	return true, nil
}

func (m *ModelSuppliers) UpdateOne(c *gin.Context) (bool, error) {
	return m.b.UpdateOne(c, m.name)
}

func (m *ModelSuppliers) UpdateMany(c *gin.Context) (bool, error) {
	// TODO
	c.JSON(400, gin.H{"message": "Not implement yet"})
	return true, nil
}

func (m *ModelSuppliers) DeleteOne(c *gin.Context) (bool, error) {
	return m.b.DeleteOne(c, m.name)
}
