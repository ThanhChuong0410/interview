package app

import (
	"errors"
	"fmt"
	"log"

	authen "github.com/chuongthanh0410/interview/authentication"
	db "github.com/chuongthanh0410/interview/database"
	_ "github.com/chuongthanh0410/interview/docs" // Đảm bảo import package docs của bạn
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	Handlers = []Basemodel{}
	Router   = &gin.Engine{}
)

// //////////////////////////////////
// Base model
// //////////////////////////////////
type Basemodel struct {
	name      string
	typeModel IBasemodel
}

type IBasemodel interface {
	Authentication(c *gin.Context) (bool, error)
	BeforeRequest(c *gin.Context) (bool, error)
	Process(c *gin.Context, function_id string) (bool, error)
	AfterRequest(c *gin.Context) (bool, error)
	ScanData(table string, getFiels []string, rows pgx.Rows) (interface{}, error)
}

func (b *Basemodel) ReturnErr(c *gin.Context, msg string) (bool, error) {
	c.JSON(
		500,
		gin.H{
			"message": msg,
			"status":  false,
			"data":    nil,
		},
	)
	return true, nil
}

func (b *Basemodel) ReturnFailed(c *gin.Context, msg string, status_code int) (bool, error) {
	c.JSON(
		status_code,
		gin.H{
			"message": msg,
			"status":  false,
			"data":    nil,
		},
	)
	return true, nil
}

func (b *Basemodel) ReturnSuc(c *gin.Context, data, metadata any) (bool, error) {
	c.JSON(
		200,
		gin.H{
			"message":  "success",
			"status":   true,
			"data":     data,
			"metadata": metadata,
		},
	)
	return true, nil
}

func (b *Basemodel) UpdateOne(c *gin.Context, table string) (bool, error) {
	id := c.Param("id")

	updateData := map[string]interface{}{}
	if err := c.BindJSON(&updateData); err != nil {
		return b.ReturnFailed(c, fmt.Sprintf("Invalid Body %s", err.Error()), 400)
	}
	setCmd, args := db.BuildSetCommand(updateData)
	err := db.DBClient.UpdateOne(table, id, setCmd, args)

	if err != nil {
		log.Printf(err.Error())
		return b.ReturnErr(c, fmt.Sprintf("update record %s error", id))
	}

	return b.ReturnSuc(c, nil, nil)
}

func (b *Basemodel) DeleteOne(c *gin.Context, table string) (bool, error) {
	id := c.Param("id")

	updateData := map[string]interface{}{}
	if err := c.BindJSON(&updateData); err != nil {
		return b.ReturnFailed(c, fmt.Sprintf("Invalid Body %s", err.Error()), 400)
	}
	err := db.DBClient.DeleteOne(table, id)

	if err != nil {
		log.Printf(err.Error())
		return b.ReturnErr(c, fmt.Sprintf("update record %s error", id))
	}

	return b.ReturnSuc(c, nil, nil)
}

func (b *Basemodel) Authentication(c *gin.Context) (bool, error) {
	fmt.Println("[base-model] authentication ...", c.ClientIP())
	if !authen.Authenticate(c) {
		return false, errors.New("authen failed")
	}
	return true, nil
}

func (b *Basemodel) BeforeRequest(c *gin.Context) (bool, error) {
	fmt.Println("[base-model] beforeRequest ...", c.ClientIP())
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		msg := "id must be uuidv4"
		b.ReturnFailed(c, msg, 400)
		return false, errors.New(msg)
	}
	return true, nil
}

func (b *Basemodel) Process(c *gin.Context) (bool, error) {
	fmt.Println("[base-model] process ...", c.ClientIP())
	return true, nil

}

func (b *Basemodel) AfterRequest(c *gin.Context) (bool, error) {
	fmt.Println("[base-model] afterRequest ...", c.ClientIP())
	return true, nil

}

func (b *Basemodel) registerRouteHandle(route *gin.Engine) {
	group := route.Group("/api/" + b.name)

	group.GET("/:id", b.BaseHandle("get_id"))
	group.GET("/", b.BaseHandle("get_all"))

	group.POST("/:id", b.BaseHandle("update_id"))
	group.POST("/", b.BaseHandle("update_many"))

	group.DELETE("/:id", b.BaseHandle("delete_id"))
}

func (b *Basemodel) BaseHandle(function_id string) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("function_id: ", function_id)
		_, err := b.typeModel.Authentication(c)
		if err != nil {
			return
		}
		_, err = b.typeModel.BeforeRequest(c)
		if err != nil {
			return
		}

		_, err = b.typeModel.Process(c, function_id)
		if err != nil {
			return
		}

		_, err = b.typeModel.AfterRequest(c)
		if err != nil {
			return
		}
	}
}

var ListModel = map[string]IBasemodel{}

func GetModel(modelName string) IBasemodel {
	m, ok := ListModel[modelName]
	if !ok {
		return nil
	}
	return m
}

func RegisterModel(model IBasemodel, doctype string) {
	o := &Basemodel{
		name:      doctype,
		typeModel: model,
	}
	ListModel[doctype] = model
	o.registerRouteHandle(Router)
}

func init() {
	Router = gin.Default()
	// Cấu hình Swagger
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Run() {
	Router.RunTLS(":8080", "./cert/cert.crt", "./cert/cert.key")
}
