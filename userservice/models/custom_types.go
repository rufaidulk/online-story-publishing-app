package models

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//todo:: find an alternate name for this struct
type Uuid struct {
	uuidStr string
}

func (u *Uuid) UuidStr() string {
	return u.uuidStr
}

// Scan implements the sql.Scanner interface
func (u *Uuid) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	return nil
}

func (u Uuid) GormDataType() string {
	return "binary"
}

func (u Uuid) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "UNHEX(REPLACE(?,'-',''))",
		Vars: []interface{}{u.uuidStr},
	}
}

func NewUuid(uuid string) Uuid {
	return Uuid{uuidStr: uuid}
}
