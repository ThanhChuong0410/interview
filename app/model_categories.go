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
type ModelCategory struct {
	b    Basemodel
	name string
}

type Category struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func init() {
	m := &ModelCategory{
		name: "categories",
	}
	RegisterModel(m, m.name)
}

func (m *ModelCategory) Authentication(c *gin.Context) (bool, error) {
	return m.b.Authentication(c)
}

func (m *ModelCategory) BeforeRequest(c *gin.Context) (bool, error) {
	return m.b.BeforeRequest(c)
}

func (m *ModelCategory) Process(c *gin.Context, function_id string) (bool, error) {
	fmt.Println("[ModelCategory-model] process ...", c.ClientIP())
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

func (m *ModelCategory) AfterRequest(c *gin.Context) (bool, error) {
	fmt.Println("[ModelCategory-model] afterRequest auditlog ...", c.ClientIP())
	return true, nil
}

func (m *ModelCategory) GetOne(c *gin.Context) (bool, error) {
	id := c.Param("id")
	rows, err := db.DBClient.SelectOne(m.name, "id", id)

	if err != nil {
		log.Printf(err.Error())
		return m.b.ReturnErr(c, fmt.Sprintf("get category %s error", id))
	}

	defer rows.Close()
	categories, err := m.ScanData(m.name, []string{}, rows)
	if err != nil {
		return m.b.ReturnErr(c, fmt.Sprintf("get category %s error %s", id, err.Error()))

	}

	return m.b.ReturnSuc(c, categories, nil)
}

func (m *ModelCategory) ScanData(table string, getFields []string, rows pgx.Rows) (interface{}, error) {
	var categories []Category
	for rows.Next() {
		category := Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during iteration: %v", err)
		return nil, errors.New("Error iterating through rows")
	}
	return categories, nil
}

func (m *ModelCategory) GetMany(c *gin.Context) (bool, error) {
	c.JSON(200, gin.H{"message": "using custom_search instead"})
	return true, nil
}

func (m *ModelCategory) UpdateOne(c *gin.Context) (bool, error) {
	return m.b.UpdateOne(c, m.name)
}

func (m *ModelCategory) UpdateMany(c *gin.Context) (bool, error) {
	// TODO
	c.JSON(400, gin.H{"message": "Not implement yet"})
	return true, nil
}

func (m *ModelCategory) DeleteOne(c *gin.Context) (bool, error) {
	return m.b.DeleteOne(c, m.name)
}
