package es

import (
	"context"
	"fmt"
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

type esOperLogStore struct {
	es *elasticsearch.TypedClient
}

var _ store.OperLogStore = &esOperLogStore{}

func (s *esOperLogStore) InsertOperLog(ctx context.Context, operlog *v1.SysOperLog, opts *api.CreateOptions) error {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.OperLogs().InsertOperLog(ctx, operlog, opts)
}

func (s *esOperLogStore) SelectOperLogList(ctx context.Context, operlog *v1.SysOperLog, opts *api.GetOptions) ([]*v1.SysOperLog, error) {
	filters := make([]types.Query, 0, 7)

	if operlog.OperIp != "" {
		cond := "*" + operlog.OperIp + "*"
		filters = append(filters, types.Query{
			Wildcard: map[string]types.WildcardQuery{
				"oper_ip": {Wildcard: &cond},
			},
		})
	}

	if operlog.Title != "" {
		cond := "*" + operlog.Title + "*"
		filters = append(filters, types.Query{
			Wildcard: map[string]types.WildcardQuery{
				"title": {Wildcard: &cond},
			},
		})
	}

	if operlog.BusinessType != nil {
		cond := operlog.BusinessType
		filters = append(filters, types.Query{
			Term: map[string]types.TermQuery{
				"business_type": {Value: &cond},
			},
		})
	}

	if len(operlog.BusinessTypes) > 0 {
		filters = append(filters, types.Query{
			Terms: &types.TermsQuery{
				TermsQuery: map[string]types.TermsQueryField{
					"business_type": operlog.BusinessTypes,
				},
			},
		})
	}

	if operlog.Status != "" {
		cond := operlog.Status
		filters = append(filters, types.Query{
			Term: map[string]types.TermQuery{
				"status": {Value: &cond},
			},
		})
	}

	if operlog.OperName != "" {
		cond := operlog.OperName
		filters = append(filters, types.Query{
			Wildcard: map[string]types.WildcardQuery{
				"oper_name": {Wildcard: &cond},
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

	resp, err := s.es.Search().Index("sys_oper_log").Query(&types.Query{
		Bool: &types.BoolQuery{Filter: filters},
	}).Sort(types.SortOptions{SortOptions: map[string]types.FieldSort{"oper_id": {Order: &sortorder.SortOrder{"desc"}}}}).Do(ctx)
	if err != nil {
		return []*v1.SysOperLog{}, err
	}
	result := make([]*v1.SysOperLog, 0, resp.Hits.Total.Value)
	for _, hit := range resp.Hits.Hits {
		var record v1.SysOperLog
		jsoniter.Unmarshal(hit.Source_, &record)
		result = append(result, &record)
	}
	return result, nil
}

func (s *esOperLogStore) DeleteOperLogByIds(ctx context.Context, operids []int64, opts *api.DeleteOptions) error {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.OperLogs().DeleteOperLogByIds(ctx, operids, opts)
}

func (s *esOperLogStore) SelectOperLogById(ctx context.Context, operid int64, opts *api.GetOptions) (*v1.SysOperLog, error) {
	resp, err := s.es.Search().Index("sys_oper_log").Query(&types.Query{
		Bool: &types.BoolQuery{Filter: []types.Query{
			{Term: map[string]types.TermQuery{"oper_id": {Value: operid}}},
		}},
	}).Do(ctx)

	if err != nil {
		return nil, err
	}
	if resp.Hits.Total.Value <= 0 {
		return nil, fmt.Errorf("未找到对应j")
	}
	var result v1.SysOperLog
	fmt.Printf("%s", resp.Hits.Hits[0].Source_)
	err = jsoniter.Unmarshal(resp.Hits.Hits[0].Source_, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *esOperLogStore) CleanOperLog(ctx context.Context, opts *api.DeleteOptions) error {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.OperLogs().CleanOperLog(ctx, opts)
}
