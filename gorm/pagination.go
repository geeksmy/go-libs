package gorm

import (
	"fmt"
	"math"

	"gorm.io/gorm"
)

// Page 分页返回
type Page struct {
	TotalRecord int64       `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	Records     interface{} `json:"records"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	NextCursor  int         `json:"next_cursor"`
}

var ErrPageParamError = fmt.Errorf("param 'result' err：%s", "unsupported destination, should be slice or struct")

// Pagination 分页
//  db 数据库连接
//  limit 每页条数
//  offset 偏移量
//  result 需要查询的结果集
func Pagination(db *gorm.DB, limit, offset int, result interface{}) (*Page, error) {
	// 如果每页条数<=0,则初始化为10条
	if limit <= 0 {
		limit = 10
	}
	// 如果偏移量小于0，则从0开始
	if offset < 0 {
		offset = 0
	}

	if result == nil {
		return nil, ErrPageParamError
	}

	var (
		page  = Page{}
		count int64
	)

	err := db.Model(result).Count(&count).Error
	if err != nil {
		return nil, err
	}

	err = db.Limit(limit).Offset(offset).Find(result).Error

	if err != nil {
		return nil, err
	}

	page.TotalRecord = count
	page.Records = result

	page.Offset = offset
	page.Limit = limit
	page.TotalPage = int(math.Ceil(float64(count) / float64(limit)))

	if count > int64(limit+offset) {
		nextCursor := offset + limit
		page.NextCursor = nextCursor
	} else {
		page.NextCursor = 0
	}

	return &page, nil
}

func PaginationJoin(db *gorm.DB, limit, offset int, model, result interface{}) (*Page, error) {
	// 如果每页条数<=0,则初始化为10条
	if limit <= 0 {
		limit = 10
	}
	// 如果偏移量小于0，则从0开始
	if offset < 0 {
		offset = 0
	}

	if result == nil {
		return nil, ErrPageParamError
	}

	var (
		page  = Page{}
		count int64
	)

	err := db.Model(model).Count(&count).Error
	if err != nil {
		return nil, err
	}

	err = db.Limit(limit).Offset(offset).Find(result).Error

	if err != nil {
		return nil, err
	}

	page.TotalRecord = count
	page.Records = result

	page.Offset = offset
	page.Limit = limit
	page.TotalPage = int(math.Ceil(float64(count) / float64(limit)))

	if count > int64(limit+offset) {
		nextCursor := offset + limit
		page.NextCursor = nextCursor
	} else {
		page.NextCursor = 0
	}

	return &page, nil
}
