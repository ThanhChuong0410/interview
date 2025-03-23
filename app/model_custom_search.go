package app

import (
	"fmt"
	"log"
	"time"

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
	return m.b.BeforeRequest(c)
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

func GetColumType(tableName string) map[string]string {
	rows, err := db.DBClient.GetInfoSchema(tableName)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

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

func (m *ModelCustomSearch) makeRowVal(
	getFields []string,
	columnDataTypes map[string]string,
) []interface{} {
	var (
		values []interface{}
		v      interface{}
	)

	for i := 0; i < len(getFields); i++ {
		columnType, e := columnDataTypes[getFields[i]]
		if !e {
			continue
		}
		if columnType == "uuid" {
			var t uuid.UUID
			values = append(values, &t)
			continue
		} else if columnType == "date" {
			var t time.Time
			values = append(values, &t)
			continue
		}
		values = append(values, &v)
	}

	return values
}

func (m *ModelCustomSearch) ScanData(table string, getFields []string, rows pgx.Rows) (interface{}, error) {
	var (
		cols            = rows.FieldDescriptions()
		rowData         = make(map[string]interface{})
		columnDataTypes = ColumnType[table]
		values          = []interface{}{}
		result          = []map[string]interface{}{}
	)
	fmt.Printf("columnDataTypes: %v\n", columnDataTypes)
	values = m.makeRowVal(getFields, columnDataTypes)
	fmt.Printf("values: %v - %v\n", values, len(values))

	if len(getFields) != 0 && len(values) != len(getFields) {
		return nil, fmt.Errorf("some field do not exist")
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			log.Printf("Error scanning table %s row: %v\n", table, err)
			continue
		}

		for i, column := range cols {
			columnType := columnDataTypes[string(column.Name)]
			if columnType == "uuid" {
				rowData[string(column.Name)] = *(values[i].(*uuid.UUID))
			} else if columnType == "date" {
				dateValue := values[i].(*time.Time)
				rowData[string(column.Name)] = dateValue.Format("2006-01-02")
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

	data, err := CustomSearch(queryCondition)
	if err != nil {
		return m.b.ReturnFailed(c, fmt.Sprintf("custom search error %s", err.Error()), 400)
	}

	metadata := gin.H{
		"page":  queryCondition.Page,
		"limit": queryCondition.Limit,
	}
	return m.b.ReturnSuc(c, data, metadata)
}

func CustomSearch(queryCondition QueryCondition) (interface{}, error) {
	var (
		data interface{}
	)
	m := GetModel("custom_search")

	if len(queryCondition.GetFields) != 0 {
		_, ok := ColumnType[queryCondition.Table]
		if !ok {
			ColumnType[queryCondition.Table] = GetColumType(queryCondition.Table)
		}
	}

	queryString, args := buildSearchQuery(queryCondition)
	rows, err := db.DBClient.Search(queryCondition.Table, queryString, args)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("get %s error: %v", queryCondition.Table, err.Error())
	}
	defer rows.Close()

	if len(queryCondition.GetFields) == 0 {
		model := GetModel(queryCondition.Table)
		data, err = model.ScanData(queryCondition.Table, []string{}, rows)
		if err != nil {
			return nil, fmt.Errorf("get %s error: %s", queryCondition.Table, err.Error())
		}
	} else {
		data, _ = m.ScanData(queryCondition.Table, queryCondition.GetFields, rows)
	}
	return data, nil
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
