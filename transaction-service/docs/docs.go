// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/transaction": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create transaction with post id, amount, etc.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transaction"
                ],
                "summary": "Create a new Transaction.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Transaction created details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TransactionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Transaction created successfully",
                        "schema": {
                            "$ref": "#/definitions/model.TransactionResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "httputil.HTTPError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "model.TransactionRequest": {
            "type": "object",
            "properties": {
                "account_name": {
                    "type": "string"
                },
                "account_number": {
                    "type": "string"
                },
                "amount": {
                    "type": "number"
                },
                "post_id": {
                    "type": "string"
                }
            }
        },
        "model.TransactionResponse": {
            "type": "object",
            "properties": {
                "account_name": {
                    "type": "string"
                },
                "account_number": {
                    "type": "string"
                },
                "amount": {
                    "type": "number"
                },
                "payment_id": {
                    "type": "string"
                },
                "payment_status": {
                    "type": "string"
                },
                "payment_url": {
                    "type": "string"
                },
                "post_id": {
                    "type": "string"
                },
                "transaction_id": {
                    "type": "string"
                },
                "user_email": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
