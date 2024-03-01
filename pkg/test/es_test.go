package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/pkg/db/doc"
	"testing"
	"time"
)

var (
	es *elasticsearch.TypedClient
)

func Init() {
	esConfig := &doc.ESConfig{
		Addrs:    []string{"https://localhost:9200"},
		Username: "sophie",
		Password: "123456",
		MaxIdle:  10,
		UseSSL:   false,
		Timeout:  5 * time.Second,
	}
	esCli, err := doc.NewES(esConfig)
	if err != nil {
		panic(err)
	}

	// es 连接测试
	ok, err := esCli.Ping().Do(context.Background())
	if err != nil || !ok {
		panic(fmt.Errorf("es 连接失败"))
	}

	es = esCli
}

func TestESSysLogininforQuery(t *testing.T) {
	cons1 := "*"
	cons2 := "*test*"
	cons3 := "2024-02-22T08:15:00"
	cons4 := "now"
	direction := &sortorder.SortOrder{"desc"}
	sorts := []types.SortCombinations{
		types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"info_id": {Order: direction},
			},
		},
	}

	resp, err := es.Search().Index("sys_logininfor").Query(&types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				{Wildcard: map[string]types.WildcardQuery{
					"ipaddr": {Wildcard: &cons1},
				}},
				{Wildcard: map[string]types.WildcardQuery{
					"user_name": {Wildcard: &cons2},
				}},
				{Term: map[string]types.TermQuery{
					"status": {Value: "0"},
				}},
				{Range: map[string]types.RangeQuery{
					"access_time": types.DateRangeQuery{
						Gte: &cons3,
						Lte: &cons4,
					},
				}},
			},
		},
	}).Sort(sorts...).Do(context.Background())
	if err != nil {
		t.Error(err)
	}

	for _, hit := range resp.Hits.Hits {
		var record v1.SysLogininfor
		err = json.Unmarshal(hit.Source_, &record)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v", record)
	}
}

func TestESSub(t *testing.T) {
	Init()

	t.Run("test-SysLogininforQuery", TestESSysLogininforQuery)
}
