package app

import (
	"fmt"
	"log"

	db "github.com/chuongthanh0410/interview/database"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

// /////////////////
// ModelCustomSearch
// /////////////////
type ModelCustomSearch struct {
	b    Basemodel
	name string
}

type QueryCondition struct {
	Filters   map[string]interface{} `json:"filters"`
	GetFields []string               `json:"get_fields"`
	Table     string                 `json:"table"`
	Limit     int                    `json:"limit"`
	Page      int                    `json:"page"`
}

var ColumnType map[string]map[string]string

func init() {
	m := &ModelCustomSearch{
		name: "custom_search",
	}
	ColumnType = make(map[string]map[string]string)
	RegisterModel(m, m.name)
}

func (m *ModelCustomSearch) Authentication(c *gin.Context) (bool, error) {
	return m.b.Authentication(c)
}

func (m *ModelCustomSearch) BeforeRequest(c *gin.Context) (bool, error) {
	fmt.Println("[ModelCustomSearch-model] beforeRequest ...", c.ClientIP())
	return true, nil
}

func (m *ModelCustomSearch) Process(c *gin.Context, function_id string) (bool, error) {
	fmt.Println("[ModelProduct-model] process ...", c.ClientIP())
	switch function_id {
	case "get_id":
		c.JSON(200, gin.H{"message": "Unknown handler!!"})

	case "update_id":
		c.JSON(200, gin.H{"message": "Unknown handler!!"})

	case "get_all":
		c.JSON(200, gin.H{"message": "Unknown handler!!"})

	case "update_many":
		return m.CustomSearch(c)

	default:
		c.JSON(200, gin.H{"message": "Unknown handler!!"})
	}
	return true, nil
}

func (m *ModelCustomSearch) AfterRequest(c *gin.Context) (bool, error) {
	fmt.Println("[ModelCustomSearch-model] afterRequest auditlog ...", c.ClientIP())
	return true, nil
}

func (m *ModelCustomSearch) GetColumType(tableName string) map[string]string {
	rows, err := db.DBClient.GetInfoSchema(tableName)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	// Map chứa kiểu dữ liệu của các cột
	columnDataTypes := make(map[string]string)
	for rows.Next() {
		var columnName, dataType string
		err := rows.Scan(&columnName, &dataType)
		if err != nil {
			log.Fatalf("Error scanning column info: %v\n", err)
		}
		columnDataTypes[columnName] = dataType
	}

	return columnDataTypes
}

func (m *ModelCustomSearch) ScanData(table string, getFields []string, rows pgx.Rows) (interface{}, error) {
	var (
		cols            = rows.FieldDescriptions()
		rowData         = make(map[string]interface{})
		columnDataTypes = ColumnType[table]
	)

	var result = []map[string]interface{}{}
	var values []interface{}
	for i := 0; i < len(getFields); i++ {
		var v interface{}
		columnType, e := columnDataTypes[getFields[i]]
		if !e {
			return nil, fmt.Errorf("field %s do not exist", getFields[i])
		}
		if columnType == "uuid" {
			var v uuid.UUID
			values = append(values, &v)
			continue
		}
		values = append(values, &v)
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}

		for i, column := range cols {
			columnType := columnDataTypes[string(column.Name)]
			if columnType == "uuid" {
				rowData[string(column.Name)] = *(values[i].(*uuid.UUID))
			} else {
				rowData[string(column.Name)] = *(values[i].(*interface{}))
			}
		}
		result = append(result, rowData)
	}
	return result, nil
}

func (m *ModelCustomSearch) CustomSearch(c *gin.Context) (bool, error) {
	var (
		data           interface{}
		queryCondition QueryCondition
	)

	if err := c.BindJSON(&queryCondition); err != nil {
		return m.b.ReturnFailed(c, fmt.Sprintf("Invalid Body %s", err.Error()), 400)
	}

	if len(queryCondition.GetFields) != 0 {
		_, ok := ColumnType[queryCondition.Table]
		if !ok {
			ColumnType[queryCondition.Table] = m.GetColumType(queryCondition.Table)
		}
	}

	queryString, args := buildSearchQuery(queryCondition)
	rows, err := db.DBClient.Search(queryCondition.Table, queryString, args)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return m.b.ReturnErr(c, fmt.Sprintf("get product error: %v", err.Error()))
	}
	defer rows.Close()

	if len(queryCondition.GetFields) == 0 {
		model := GetModel(queryCondition.Table)
		data, err = model.ScanData(queryCondition.Table, []string{}, rows)
		if err != nil {
			return m.b.ReturnFailed(c, fmt.Sprintf("get data error: %s", err.Error()), 400)
		}
	} else {
		data, _ = m.ScanData(queryCondition.Table, queryCondition.GetFields, rows)
	}

	metadata := gin.H{
		"page":  queryCondition.Page,
		"limit": queryCondition.Limit,
	}
	return m.b.ReturnSuc(c, data, metadata)
}

func buildDynamicSelectQuery(table string, fields []string) string {
	if len(fields) == 0 {
		return fmt.Sprintf("SELECT * FROM %s WHERE 1=1", table)
	}

	fieldList := fmt.Sprintf("%s", fields[0])
	for _, field := range fields[1:] {
		fieldList = fmt.Sprintf("%s, %s", fieldList, field)
	}

	return fmt.Sprintf("SELECT %s FROM %s WHERE 1=1", fieldList, table)
}

func buildSearchQuery(queryCondition QueryCondition) (string, []interface{}) {
	var (
		filters   = queryCondition.Filters
		getFields = queryCondition.GetFields
		limit     = queryCondition.Limit
		offset    = queryCondition.Page
		table     = queryCondition.Table
		// sorts     = queryCondition.Sorts
	)

	query := buildDynamicSelectQuery(table, getFields)
	var args []interface{}
	var argIndex int

	for k, v := range filters {
		argIndex++
		query += fmt.Sprintf(" AND %s = $%d", k, argIndex)
		args = append(args, v)
	}

	if queryCondition.Page == 0 {
		queryCondition.Page = 1
	}
	offset = (queryCondition.Page - 1) * queryCondition.Limit

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// TODO sort not implement

	return query, args
}
