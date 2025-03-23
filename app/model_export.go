package app

import (
	"fmt"
	"log"
	"reflect"
	"time"

	db "github.com/chuongthanh0410/interview/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jung-kurt/gofpdf"
)

// /////////////////
// modelA request
// /////////////////
type ModelExports struct {
	b    Basemodel
	name string
}

func init() {
	m := &ModelExports{
		name: "exports",
	}
	RegisterModel(m, m.name)
}

func (m *ModelExports) Authentication(c *gin.Context) (bool, error) {
	return m.b.Authentication(c)
}

func (m *ModelExports) BeforeRequest(c *gin.Context) (bool, error) {
	return m.b.BeforeRequest(c)
}

func (m *ModelExports) Process(c *gin.Context, function_id string) (bool, error) {
	fmt.Println("[ModelExports-model] process ...", c.ClientIP())
	switch function_id {
	case "get_id":
		return m.GetOne(c)

	case "update_id":
		return m.UpdateOne(c)

	case "get_all":
		return m.GetMany(c)

	case "update_many":
		return m.Exports(c)

	case "delete_id":
		return m.DeleteOne(c)

	default:
		c.JSON(200, gin.H{"message": "Unknown handler!!"})
	}
	return true, nil
}

func (m *ModelExports) AfterRequest(c *gin.Context) (bool, error) {
	fmt.Println("[ModelExports-model] afterRequest auditlog ...", c.ClientIP())
	return true, nil
}

func (m *ModelExports) GetOne(c *gin.Context) (bool, error) {
	// TODO
	c.JSON(400, gin.H{"message": "Not implement yet"})
	return true, nil
}

func (m *ModelExports) ScanData(table string, getFields []string, rows pgx.Rows) (interface{}, error) {
	return nil, nil
}

func (m *ModelExports) GetMany(c *gin.Context) (bool, error) {
	// TODO
	c.JSON(400, gin.H{"message": "Not implement yet"})
	return true, nil
}

func (m *ModelExports) UpdateOne(c *gin.Context) (bool, error) {
	// TODO
	c.JSON(400, gin.H{"message": "Not implement yet"})
	return true, nil
}

func (m *ModelExports) Exports(c *gin.Context) (bool, error) {
	var (
		data           interface{}
		queryCondition QueryCondition
		model          IBasemodel
	)

	if err := c.BindJSON(&queryCondition); err != nil {
		return m.b.ReturnFailed(c, fmt.Sprintf("Invalid Body %s", err.Error()), 400)
	}

	_, ok := ColumnType[queryCondition.Table]
	if !ok {
		ColumnType[queryCondition.Table] = GetColumType(queryCondition.Table)
	}

	queryString, args := buildSearchQuery(queryCondition)
	rows, err := db.DBClient.Search(queryCondition.Table, queryString, args)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return m.b.ReturnFailed(c,
			fmt.Sprintf("search %s error: %v", queryCondition.Table, err.Error()),
			400)
	}
	defer rows.Close()

	if len(queryCondition.GetFields) == 0 {
		model = GetModel(queryCondition.Table)

	} else {
		model = GetModel("custom_search")
	}
	data, err = model.ScanData(queryCondition.Table, queryCondition.GetFields, rows)
	if err != nil {
		log.Printf("Error scan data: %v", err)
		return m.b.ReturnFailed(c,
			fmt.Sprintf("scan %s error: %v", queryCondition.Table, err.Error()),
			400)
	}

	_, err = ExportPDF(data)
	if err != nil {
		return m.b.ReturnFailed(c, fmt.Sprintf("custom search error %s", err.Error()), 400)
	}

	c.JSON(400, gin.H{"message": "Not implement yet"})
	return true, nil
}

func exportPDFaMap(input interface{}) (bool, error) {
	//
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	value := reflect.ValueOf(input)
	if value.Kind() != reflect.Slice || value.Len() == 0 {
		return false, fmt.Errorf("invalid data export or empty")
	}

	firstItem := value.Index(0).Interface()
	var headers []string

	if record, ok := firstItem.(map[string]interface{}); ok {
		for key := range record {
			headers = append(headers, key)
		}
	}

	for _, header := range headers {
		pdf.Cell(40, 10, header)
	}
	pdf.Ln(10)

	for i := 0; i < value.Len(); i++ {
		record := value.Index(i).Interface().(map[string]interface{})
		for _, header := range headers {
			pdf.Cell(40, 10, fmt.Sprintf("%v", record[header]))
		}
		pdf.Ln(10)
	}

	// Save
	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		return false, err
	}
	return true, nil
}

func exportPDFaStruct(input any) (bool, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	value := reflect.ValueOf(input)

	if value.Kind() != reflect.Slice || value.Len() == 0 {
		return false, fmt.Errorf("export data must not empty")
	}

	firstItem := value.Index(0)
	if firstItem.Kind() != reflect.Struct {
		return false, fmt.Errorf("export data must be slice struct")
	}

	var headers []string
	var fields []reflect.StructField

	for i := 0; i < firstItem.NumField(); i++ {
		field := firstItem.Type().Field(i)
		headers = append(headers, field.Name)
		fields = append(fields, field)
	}

	for _, header := range headers {
		pdf.Cell(40, 10, header)
	}
	pdf.Ln(10)

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		for _, field := range fields {
			if field.Type.Name() == "Time" {
				pdf.MultiCell(40, 10,
					fmt.Sprintf("%v",
						item.FieldByName(field.Name).Interface().(time.Time).Format("2006-01-02")),
					"",
					"L",
					false,
				)
				continue
			}
			pdf.MultiCell(40, 10,
				fmt.Sprintf("%v", item.FieldByName(field.Name).Interface()),
				"",
				"L",
				false,
			)
		}
		pdf.Ln(10)
	}
	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		log.Fatal(err)
	}

	return true, nil
}

func ExportPDF(input interface{}) (bool, error) {
	_, ok := input.([]map[string]interface{})
	if !ok {
		return exportPDFaStruct(input)
	}
	return exportPDFaMap(input)
}

func (m *ModelExports) DeleteOne(c *gin.Context) (bool, error) {
	// TODO
	c.JSON(400, gin.H{"message": "Not implement yet"})
	return true, nil
}
