package gormutil

import (
	"gorm.io/gorm"
)

// Scope func(db *gorm.DB) *gorm.DB 的 type alias
// example:
// 		myScope Scope = func(db *gorm.DB) {return db}
// 		db.Scopes(myScope)
type Scope = func(db *gorm.DB) *gorm.DB
