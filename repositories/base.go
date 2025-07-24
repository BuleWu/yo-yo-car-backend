package repositories

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type SearchQuery struct {
	Column   string      `json:"column"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}
