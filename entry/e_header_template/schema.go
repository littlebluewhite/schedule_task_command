package e_header_template

import "github.com/goccy/go-json"

type HeaderTemplate struct {
	ID   int32           `json:"id"`
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type HeaderTemplateCreate struct {
	Name string          `json:"name" binding:"required"`
	Data json.RawMessage `json:"data" binding:"required"`
}

type HeaderTemplateUpdate struct {
	ID   int32            `json:"id" binding:"required"`
	Name *string          `json:"name"`
	Data *json.RawMessage `json:"data"`
}
