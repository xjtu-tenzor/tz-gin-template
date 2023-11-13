package model

import (
	"template/common"

	"gorm.io/gorm"
)

func Paginate(pager common.PagerForm) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pager.Page <= 0 {
			pager.Page = 1
		}

		switch {
		case pager.Limit > 20:
			pager.Limit = 20
		case pager.Limit <= 0:
			pager.Limit = 10
		}

		offset := (pager.Page - 1) * pager.Limit
		return db.Offset(offset).Limit(pager.Limit)
	}
}
