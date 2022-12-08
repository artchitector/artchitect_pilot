package model

import (
	"database/sql"
	"gorm.io/gorm"
)

// TODO Need to move postgres-model to separate package and use it in both services
type Painting struct {
	gorm.Model
	Caption string
	Bytes   sql.RawBytes
}
