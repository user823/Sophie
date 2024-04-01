package utils

import (
	"github.com/user823/Sophie/api"
	"gorm.io/gorm"
)

func CountQuery(query *gorm.DB, opts *api.GetOptions, column string) (res int64) {
	realQuery := opts.SQLConditionWithoutPage(query, column)
	realQuery.Count(&res)
	return
}
