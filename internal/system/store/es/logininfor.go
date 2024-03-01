package es

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/store/mysql"
	"github.com/user823/Sophie/pkg/utils"
)

type esLogininforStore struct {
	es *elasticsearch.TypedClient
}

var _ store.LogininforStore = &esLogininforStore{}

func (s *esLogininforStore) InsertLogininfor(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.CreateOptions) error {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Logininfors().InsertLogininfor(ctx, logininfor, opts)
}

func (s *esLogininforStore) SelectLogininforList(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.GetOptions) ([]*v1.SysLogininfor, error) {
	filters := make([]types.Query, 0, 4)
	if logininfor.Ipaddr != "" {
		cond := "*" + logininfor.Ipaddr + "*"
		filters = append(filters, types.Query{
			Wildcard: map[string]types.WildcardQuery{
				"ip_addr": {Wildcard: &cond},
			},
		})
	}

	if logininfor.Status != "" {
		cond := logininfor.Status
		filters = append(filters, types.Query{
			Term: map[string]types.TermQuery{
				"status": {Value: &cond},
			},
		})
	}

	if logininfor.UserName != "" {
		cond := "*" + logininfor.UserName + "*"
		filters = append(filters, types.Query{
			Wildcard: map[string]types.WildcardQuery{
				"user_name": {Wildcard: &cond},
			},
		})
	}

	if opts.BeginTime != 0 || opts.EndTime != 0 {
		tr := types.NewDateRangeQuery()
		if opts.BeginTime != 0 {
			beginTime := utils.Second2Time(opts.BeginTime).String()
			tr.Gte = &beginTime
		}
		if opts.EndTime != 0 {
			endTime := utils.Second2Time(opts.EndTime).String()
			tr.Lte = &endTime
		}
		filters = append(filters, types.Query{
			Range: map[string]types.RangeQuery{
				"access_time": tr,
			},
		})
	}

	resp, err := s.es.Search().Index("sys_logininfor").Query(&types.Query{
		Bool: &types.BoolQuery{Filter: filters},
	}).Sort(types.SortOptions{SortOptions: map[string]types.FieldSort{"info_id": {Order: &sortorder.SortOrder{"desc"}}}}).Do(ctx)

	if err != nil {
		return []*v1.SysLogininfor{}, err
	}

	result := make([]*v1.SysLogininfor, 0, resp.Hits.Total.Value)
	for _, hit := range resp.Hits.Hits {
		var record v1.SysLogininfor
		jsoniter.Unmarshal(hit.Source_, &record)
		result = append(result, &record)
	}
	return result, nil
}

func (s *esLogininforStore) DeleteLogininforByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	//_, err := s.es.DeleteByQuery("sys_logininfor").Query(&types.Query{
	//	Terms: &types.TermsQuery{
	//		TermsQuery: map[string]types.TermsQueryField{
	//			"info_id": ids,
	//		},
	//	},
	//}).Do(ctx)
	//if err != nil {
	//	return err
	//}

	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Logininfors().DeleteLogininforByIds(ctx, ids, opts)
}

func (s *esLogininforStore) CleanLogininfor(ctx context.Context, opts *api.DeleteOptions) error {
	//_, err := s.es.DeleteByQuery("sys_logininfor").Query(&types.Query{
	//	MatchAll: &types.MatchAllQuery{},
	//}).Do(ctx)
	//if err != nil {
	//	return err
	//}

	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Logininfors().CleanLogininfor(ctx, opts)
}
