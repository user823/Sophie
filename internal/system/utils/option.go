package utils

import (
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
)

func BuildGetOption(pageInfo *v1.PageInfo, daterange *v1.DateRange, useCache bool) *api.GetOptions {
	opts := &api.GetOptions{
		Cache: useCache,
	}

	if pageInfo != nil {
		opts.PageNum = pageInfo.GetPageNum()
		opts.PageSize = pageInfo.GetPageSize()
		opts.OrderByColumn = pageInfo.GetOrderByColumn()
		opts.IsAsc = pageInfo.GetIsAsc() == api.IS_ASC
	}

	if daterange != nil {
		opts.BeginTime = daterange.BeginTime
		opts.EndTime = daterange.EndTime
	}

	return opts
}
