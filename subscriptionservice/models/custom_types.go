package models

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UuidData struct {
	uuidStr string
}

func (u *UuidData) UuidStr() string {
	return u.uuidStr
}

// Scan implements the sql.Scanner interface
func (u *UuidData) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	return nil
}

func (u UuidData) GormDataType() string {
	return "binary"
}

func (u UuidData) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "UNHEX(REPLACE(?,'-',''))",
		Vars: []interface{}{u.uuidStr},
	}
}

func NewUuidData(uuid string) UuidData {
	return UuidData{uuidStr: uuid}
}
