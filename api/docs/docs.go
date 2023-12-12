// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/media/photo": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Through this api frontent can upload photo and get the link to the media.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "Upload media",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.MediaResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/template": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Here template can be updated.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "Update template",
                "parameters": [
                    {
                        "description": "post info",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TemplateUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.TemplateApiResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Here template can be created.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "Create template",
                "parameters": [
                    {
                        "description": "post info",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TemplateCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.TemplateApiResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/template/list": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Here all templates can be got.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "Get templates list",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "order_by_created_at",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.TemplateApiFindResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/template/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Here template can be got.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "Get template by key",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.TemplateApiResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Here template can be deleted.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "Delete template",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/user": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Here user can be updated.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update user",
                "parameters": [
                    {
                        "description": "post info",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserApiUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserApiResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Here user can be registered.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Authorzation"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "post info",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserRegisterReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserApiResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Here user can be deleted, user_id is taken from token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Delete user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/user/check/{email}": {
            "get": {
                "description": "Here user status is checked. If user is exists in database it should be logged in else registered",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Authorzation"
                ],
                "summary": "Check User status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "email",
                        "name": "email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserApiResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/user/forgot-password/verify": {
            "post": {
                "description": "Through this api user forgot  password can be enabled.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Authorzation"
                ],
                "summary": "User forgot password",
                "parameters": [
                    {
                        "description": "User Login",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserForgotPasswordVerifyReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/user/forgot-password/{user_name_or_email}": {
            "get": {
                "description": "Through this api user forgot  password can be enabled.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Authorzation"
                ],
                "summary": "User forgot password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user_name_or_email",
                        "name": "user_name_or_email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "Through this api user is logged in",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Authorzation"
                ],
                "summary": "User Login",
                "parameters": [
                    {
                        "description": "User Login",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserApiResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/user/otp": {
            "get": {
                "description": "Here otp can be checked if true.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Authorzation"
                ],
                "summary": "Check Otp",
                "parameters": [
                    {
                        "type": "string",
                        "description": "email",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "otp",
                        "name": "otp",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.OtpCheckResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        },
        "/user/profile": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Here user profile info can be got by id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get user by key",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserApiResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/models.StandardResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.MediaResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/models.UploadPhotoRes"
                },
                "error_code": {
                    "type": "integer"
                },
                "error_message": {
                    "type": "string"
                }
            }
        },
        "models.OtpCheckResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "object",
                    "properties": {
                        "is_right": {
                            "type": "boolean"
                        }
                    }
                },
                "error_code": {
                    "type": "integer"
                },
                "error_message": {
                    "type": "string"
                }
            }
        },
        "models.StandardResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "error_code": {
                    "type": "integer"
                },
                "error_message": {
                    "type": "string"
                }
            }
        },
        "models.TemplateApiFindResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/models.TemplateFindResponse"
                },
                "error_code": {
                    "type": "integer"
                },
                "error_message": {
                    "type": "string"
                }
            }
        },
        "models.TemplateApiResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/models.TemplateResponse"
                },
                "error_code": {
                    "type": "integer"
                },
                "error_message": {
                    "type": "string"
                }
            }
        },
        "models.TemplateCreateReq": {
            "type": "object",
            "properties": {
                "template_name": {
                    "type": "string"
                }
            }
        },
        "models.TemplateFindResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "templates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.TemplateResponse"
                    }
                }
            }
        },
        "models.TemplateResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "template_name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.TemplateUpdateReq": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "template_name": {
                    "type": "string"
                }
            }
        },
        "models.UploadPhotoRes": {
            "type": "object",
            "properties": {
                "photo_url": {
                    "type": "string"
                }
            }
        },
        "models.UserApiResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/models.UserResponse"
                },
                "error_code": {
                    "type": "integer"
                },
                "error_message": {
                    "type": "string"
                }
            }
        },
        "models.UserApiUpdateReq": {
            "type": "object",
            "properties": {
                "user_name": {
                    "type": "string"
                }
            }
        },
        "models.UserForgotPasswordVerifyReq": {
            "type": "object",
            "properties": {
                "new_password": {
                    "type": "string"
                },
                "otp": {
                    "type": "string"
                },
                "user_name_or_email": {
                    "type": "string"
                }
            }
        },
        "models.UserLoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "user_name_or_email": {
                    "type": "string"
                }
            }
        },
        "models.UserRegisterReq": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "otp": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                }
            }
        },
        "models.UserResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Monolithic project API Endpoints",
	Description:      "Here QA can test and frontend or mobile developers can get information of API endpoints.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
