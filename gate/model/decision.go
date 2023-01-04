package model

import (
	"database/sql"
	"gorm.io/gorm"
)

// Decision made with help of Origin
type Decision struct {
	gorm.Model
	Output              float64      // what value was given from Origin
	Artifact            sql.RawBytes // random seed used by number generation
	ArtifactContentType string       // type of attached artifact
}
