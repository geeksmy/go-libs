package gorm

import (
	"gorm.io/gorm"
)

// PaginationScope 分页 Scope,
//  lastID: 前一页的最后一行记录的 ID
//  perPage: 每页行数, default 100
func PaginationScope(lastID, perPage int) Scope {
	if perPage <= 0 {
		perPage = 100
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id>?", lastID).Limit(perPage)
	}
}
