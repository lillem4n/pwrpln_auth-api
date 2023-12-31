// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Power Plan",
            "url": "https://pwrpln.com/",
            "email": "lilleman@larvit.se"
        },
        "license": {
            "name": "MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/accounts": {
            "get": {
                "description": "Requires Authorization-header with role \"admin\".\nExample: Authorization: bearer xxx\nWhere \"xxx\" is a valid JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get accounts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Account"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Requires Authorization-header with role \"admin\".\nExample: Authorization: bearer xxx\nWhere \"xxx\" is a valid JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create an account",
                "operationId": "account-create",
                "parameters": [
                    {
                        "description": "Account object to be written to database",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.AccountInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/db.CreatedAccount"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    }
                }
            }
        },
        "/accounts/:id": {
            "delete": {
                "description": "Requires Authorization-header with role \"admin\" or a matching account id\nExample: Authorization: bearer xxx\nWhere \"xxx\" is a valid JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete an account",
                "operationId": "account-del",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    }
                }
            }
        },
        "/accounts/{id}": {
            "get": {
                "description": "Requires Authorization-header with either role \"admin\" or with a matching account id.\nExample: Authorization: bearer xxx\nWhere \"xxx\" is a valid JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get account by id",
                "operationId": "get-account-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Account"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    }
                }
            }
        },
        "/accounts/{id}/fields": {
            "put": {
                "description": "Requires Authorization-header with role \"admin\".\nExample: Authorization: bearer xxx\nWhere \"xxx\" is a valid JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update account fields",
                "operationId": "account-update-fields",
                "parameters": [
                    {
                        "description": "Fields array with objects to be written to database",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.AccountCreateInputFields"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Account"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    }
                }
            }
        },
        "/auth/api-key": {
            "post": {
                "description": "Authenticate account by API Key",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Authenticate account by API Key",
                "operationId": "auth-account-by-api-key",
                "parameters": [
                    {
                        "description": "API Key as a string in JSON format (just encapsulate the string with \\",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.ResToken"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    }
                }
            }
        },
        "/auth/password": {
            "post": {
                "description": "Authenticate account by Password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Authenticate account by Password",
                "operationId": "auth-account-by-password",
                "parameters": [
                    {
                        "description": "Name and password to auth by",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.AuthInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.ResToken"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    }
                }
            }
        },
        "/renew-token": {
            "post": {
                "description": "Renew token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Renew token",
                "operationId": "renew-token",
                "parameters": [
                    {
                        "description": "Renewal token as a string in JSON format (just encapsulate the string with \\",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.ResToken"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ResJSONError"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.Account": {
            "type": "object",
            "properties": {
                "created": {
                    "type": "string"
                },
                "fields": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "db.AccountCreateInputFields": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "values": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "db.CreatedAccount": {
            "type": "object",
            "properties": {
                "apiKey": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handlers.AccountInput": {
            "type": "object",
            "properties": {
                "fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/db.AccountCreateInputFields"
                    }
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "handlers.AuthInput": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "handlers.ResJSONError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "field": {
                    "type": "string"
                }
            }
        },
        "handlers.ResToken": {
            "type": "object",
            "properties": {
                "jwt": {
                    "type": "string"
                },
                "renewalToken": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "JWT Auth API",
	Description:      "This is a tiny http API for auth. Register accounts, auth with api-key or name/password, renew JWT tokens...",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
