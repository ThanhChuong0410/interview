// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/products/:id": {
            "get": {
                "description": "Get one products",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get one products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/app.Product"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.Product": {
            "type": "object",
            "properties": {
                "added_date": {
                    "description": "@Description Date when the product was added\n@example \"2025-03-22\" \"YYYY-MM-DD\"\n@format date",
                    "type": "string"
                },
                "category_id": {
                    "description": "@format uuid",
                    "type": "string"
                },
                "id": {
                    "description": "@format uuid",
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                },
                "reference": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "stock_city": {
                    "type": "string"
                },
                "supplier_id": {
                    "description": "@format uuid",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/",
	Schemes:          []string{},
	Title:            "Gin Swagger Example API",
	Description:      "API xác thực dùng token tĩnh, giá trị = 1234567890abcdefjustforspeed",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
