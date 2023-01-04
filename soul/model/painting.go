package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Painting struct {
	gorm.Model
	Caption string
	Bytes   sql.RawBytes
}
