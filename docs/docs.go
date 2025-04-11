// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/free-tables": {
            "get": {
                "description": "Возвращает список доступных столиков на указанный период",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tables"
                ],
                "summary": "Получить список свободных столиков",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Начало периода бронирования (в формате RFC3339)",
                        "name": "reservation_time_from",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Конец периода бронирования (в формате RFC3339)",
                        "name": "reservation_time_to",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.FreeTable"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/reservations": {
            "post": {
                "description": "Эта функция выполняет бронирование столика для пользователя на указанный период времени.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reservations"
                ],
                "summary": "Бронирование столика",
                "parameters": [
                    {
                        "description": "Информация о бронировании",
                        "name": "reservation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ReservationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "dto.FreeTable": {
            "type": "object",
            "properties": {
                "restaurant_name": {
                    "type": "string"
                },
                "seats_number": {
                    "type": "integer"
                },
                "table_id": {
                    "type": "integer"
                },
                "table_number": {
                    "type": "integer"
                }
            }
        },
        "dto.ReservationRequest": {
            "type": "object",
            "properties": {
                "reservation_time_from": {
                    "description": "Время начала бронирования (RFC3339)",
                    "type": "string"
                },
                "reservation_time_to": {
                    "description": "Время окончания бронирования (RFC3339)",
                    "type": "string"
                },
                "table_id": {
                    "description": "Идентификатор столика",
                    "type": "integer"
                },
                "user_id": {
                    "description": "Идентификатор пользователя",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8082",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Restorio API",
	Description:      "API для бронирования столиков",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
