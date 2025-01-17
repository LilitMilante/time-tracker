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
        "/users": {
            "get": {
                "description": "Get a list of users with optional filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of users per page",
                        "name": "per_page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Passport series",
                        "name": "passport_series",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Passport number",
                        "name": "passport_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Surname",
                        "name": "surname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Patronymic",
                        "name": "patronymic",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Address",
                        "name": "address",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/tracker.User"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new user with passport number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "Passport number in format '1234 567890'",
                        "name": "passportNumber",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tracker.PassportNumber"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tracker.User"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "patch": {
                "description": "Update user details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update an existing user",
                "parameters": [
                    {
                        "description": "User to update",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tracker.UpdateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tracker.User"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/{user_id}": {
            "delete": {
                "description": "Delete a user by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Delete a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid user ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/{user_id}/report": {
            "get": {
                "description": "Get the time spent on tasks by a user within a specified period",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Get task spend times by user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Start date in format 'DD-MM-YYYY'",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date in format 'DD-MM-YYYY'",
                        "name": "end_date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/tracker.TaskSpendTime"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "User or task not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/work/finish": {
            "post": {
                "description": "Finish work on a task for a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "work"
                ],
                "summary": "Finish work on a task",
                "parameters": [
                    {
                        "description": "Finish work request",
                        "name": "finishWorkRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tracker.FinishWorkRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Work finished",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Task not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/work/start": {
            "post": {
                "description": "Start work on a task for a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "work"
                ],
                "summary": "Start work on a task",
                "parameters": [
                    {
                        "description": "Start work request",
                        "name": "startWorkRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tracker.StartWorkRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Work started",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "tracker.FinishWorkRequest": {
            "type": "object",
            "properties": {
                "task_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "tracker.PassportNumber": {
            "type": "object",
            "properties": {
                "passportNumber": {
                    "type": "string"
                }
            }
        },
        "tracker.StartWorkRequest": {
            "type": "object",
            "properties": {
                "task_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "tracker.TaskSpendTime": {
            "type": "object",
            "properties": {
                "spend_time_sec": {
                    "type": "integer"
                },
                "task_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "tracker.UpdateUser": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "passport_number": {
                    "type": "integer"
                },
                "passport_series": {
                    "type": "integer"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "tracker.User": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "passport_number": {
                    "type": "integer"
                },
                "passport_series": {
                    "type": "integer"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
