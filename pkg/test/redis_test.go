package test

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/db/kv/redis"
	"github.com/user823/Sophie/pkg/utils"
	"reflect"
	"testing"
	"time"
)

var (
	client           kv.RedisStore
	ctx              = context.Background()
	connectionConfig = &redis.RedisConfig{
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "123456",
		Database: 0,
	}
)

const (
	testKeyPrefix string = "testing-"
	hashKey       bool   = false
	randomExp     bool   = true
)

func init() {
	go redis.KeepConnection(ctx, connectionConfig)
	client = kv.NewKVStore("redis").(kv.RedisStore)
	client.SetKeyPrefix(testKeyPrefix)
	client.SetHashKey(hashKey)
	client.SetRandomExp(randomExp)
	time.Sleep(2 * time.Second)
}

func Test_RedisClient_Connect(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx    context.Context
		config any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test1", fields: fields{"test_connection", false}, args: args{ctx, connectionConfig}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.Connect(tt.args.ctx, tt.args.config)
			fmt.Println(r.Connected())
		})
	}
}

func Test_RedisClient_SetKey(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx    context.Context
		key    string
		value  string
		expire int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "a", "b", 0}, false},
		{"test2", fields{testKeyPrefix, hashKey}, args{ctx, "key", "", utils.SecondToNano(3600)}, false},
		{"test3", fields{testKeyPrefix, hashKey}, args{ctx, "c", "d", utils.SecondToNano(3600)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.SetKey(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expire); (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.SetKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_RedisClient_GetKey(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx     context.Context
		keyname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "key"}, "", false},
		{"test2", fields{testKeyPrefix, hashKey}, args{ctx, "a"}, "b", false},
	}
	// redis 中需要先存放key: a value: b, 注意，开启hashKey的时候要先测试SetKey
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := client
			got, err := r.GetKey(tt.args.ctx, tt.args.keyname)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisClient.GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_GetMultiKey(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, []string{"a", "c"}}, []string{"b", "d"}, false},
		{"test2", fields{testKeyPrefix, hashKey}, args{ctx, []string{"e"}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetMultiKey(tt.args.ctx, tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.GetMultiKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisClient.GetMultiKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_GetKeyTTL(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx     context.Context
		keyname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "c"}, 0, false},
		{"test2", fields{testKeyPrefix, hashKey}, args{ctx, "a"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetKeyTTL(tt.args.ctx, tt.args.keyname)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.GetKeyTTL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("RedisClient.GetKeyTTL() = %v", got)
		})
	}
}

func Test_RedisClient_GetRawKey(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx     context.Context
		keyName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "a"}, "b", false},
		{"test2", fields{testKeyPrefix, hashKey}, args{ctx, "c"}, "", true},
	}
	// redis 实现存放Raw key: a, values: b
	client.SetRawKey(ctx, "a", "b", 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetRawKey(tt.args.ctx, tt.args.keyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.GetRawKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisClient.GetRawKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_GetKeys(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx    context.Context
		filter string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, ""}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.GetKeys(tt.args.ctx, tt.args.filter); got != nil {
				//t.Errorf("RedisClient.GetKeys() = %v, want %v", got, tt.want)
				t.Logf("RedisClient.GetKeys() = %v", got)
			}
		})
	}
}

func Test_RedisClient_DeleteKey(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "a"}, true},
		{"test2", fields{testKeyPrefix, hashKey}, args{ctx, "a"}, false},
	}
	// 实现要有 key: a value: b
	client.SetKey(ctx, "a", "b", 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteKey(tt.args.ctx, tt.args.key); got != tt.want {
				t.Errorf("RedisClient.DeleteKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_DeleteAllKeys(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx}, true},
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteAllKeys(tt.args.ctx); got != tt.want {
				t.Errorf("RedisClient.DeleteAllKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_DeleteRawKey(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "a"}, true},
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "a"}, false},
	}
	client.SetRawKey(ctx, "a", "b", 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteRawKey(tt.args.ctx, tt.args.key); got != tt.want {
				t.Errorf("RedisClient.DeleteRawKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_GetKeysAndValuesWithFilter(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx    context.Context
		filter string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "a"}, map[string]string{"a": "b"}},
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "c"}, map[string]string{"c": "d", "cx": "e"}},
	}
	client.SetRawKey(ctx, "a", "b", 0)
	client.SetRawKey(ctx, "c", "d", 0)
	client.SetRawKey(ctx, "cx", "e", 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.GetKeysAndValuesWithFilter(tt.args.ctx, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisClient.GetKeysAndValuesWithFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_GetKeysAndValues(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := client
			if got := r.GetKeysAndValues(tt.args.ctx); got != nil {
				t.Logf("RedisClient.GetKeysAndValues() = %v", got)
			}

		})
	}
}

func Test_RedisClient_DeleteKeys(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, []string{"a"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteKeys(tt.args.ctx, tt.args.keys); got != tt.want {
				t.Errorf("RedisClient.DeleteKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_IncrememntWithExpire(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx    context.Context
		key    string
		expire int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "cnt", 0}, 1},
		{"test2", fields{testKeyPrefix, hashKey}, args{ctx, "cnt", 0}, 2},
	}
	client.DeleteKey(ctx, "cnt")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.IncrememntWithExpire(tt.args.ctx, tt.args.key, tt.args.expire); got != tt.want {
				t.Errorf("RedisClient.IncrememntWithExpire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_Decrement(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wants  int64
	}{
		{"test1", fields{testKeyPrefix, hashKey}, args{ctx, "cnt"}, 4},
		{"test2", fields{testKeyPrefix, hashKey}, args{ctx, "cnt"}, 3},
	}
	client.DeleteKey(ctx, "cnt")
	client.SetKey(ctx, "cnt", "5", 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := client
			if got := r.Decrement(tt.args.ctx, tt.args.key); got != tt.wants {
				t.Errorf("RedisClient.Decrement = %v, want %v", got, tt.wants)
			}
		})
	}
}

func Test_RedisClient_SetRollingWindow(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx      context.Context
		key      string
		per      int64
		val      string
		pipeline bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, got1 := r.SetRollingWindow(tt.args.ctx, tt.args.key, tt.args.per, tt.args.val, tt.args.pipeline)
			if got != tt.want {
				t.Errorf("RedisClient.SetRollingWindow() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RedisClient.SetRollingWindow() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_RedisClient_GetRollingWindow(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx      context.Context
		key      string
		per      int64
		pipeline bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, got1 := r.GetRollingWindow(tt.args.ctx, tt.args.key, tt.args.per, tt.args.pipeline)
			if got != tt.want {
				t.Errorf("RedisClient.GetRollingWindow() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RedisClient.GetRollingWindow() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_RedisClient_GetSet(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetSet(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.GetSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisClient.GetSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_AddToSet(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.AddToSet(tt.args.ctx, tt.args.key, tt.args.value)
		})
	}
}

func Test_RedisClient_GetAndDeleteSet(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.GetAndDeleteSet(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisClient.GetAndDeleteSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_RemoveFromSet(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.RemoveFromSet(tt.args.ctx, tt.args.key, tt.args.value)
		})
	}
}

func Test_RedisClient_DeleteScanMatch(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx    context.Context
		filter string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteScanMatch(tt.args.ctx, tt.args.filter); got != tt.want {
				t.Errorf("RedisClient.DeleteScanMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_AddToSortedSet(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
		score float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.AddToSortedSet(tt.args.ctx, tt.args.key, tt.args.value, tt.args.score)
		})
	}
}

func Test_RedisClient_GetSortedSetRange(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx  context.Context
		key  string
		from string
		to   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		want1   []float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, got1, err := r.GetSortedSetRange(tt.args.ctx, tt.args.key, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.GetSortedSetRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisClient.GetSortedSetRange() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RedisClient.GetSortedSetRange() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_RedisClient_RemoveSortedSetRange(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx  context.Context
		key  string
		from string
		to   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.RemoveSortedSetRange(tt.args.ctx, tt.args.key, tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.RemoveSortedSetRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_RedisClient_GetListRange(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx  context.Context
		key  string
		from int64
		to   int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetListRange(tt.args.ctx, tt.args.key, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.GetListRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisClient.GetListRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_RemoveFromList(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.RemoveFromList(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.RemoveFromList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_RedisClient_AppendToSet(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.AppendToSet(tt.args.ctx, tt.args.key, tt.args.value)
		})
	}
}

func Test_RedisClient_Exists(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.Exists(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisClient.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisClient_SetExp(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx    context.Context
		key    string
		expire int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.SetExp(tt.args.ctx, tt.args.key, tt.args.expire); (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.SetExp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_RedisClient_SetRawKey(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ctx    context.Context
		key    string
		value  string
		expire int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redis.RedisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.SetRawKey(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expire); (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.SetRawKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisAppendSetPipelined(t *testing.T) {
	client.AppendToSetPipelined(ctx, "test", []string{"123", "456"})
	result := client.GetAndDeleteSet(ctx, "test")
	fmt.Println(result)
}

func TestGetAndDelete(t *testing.T) {
	if err := client.SetKey(ctx, "tbb", "b", utils.SecondToNano(600)); err != nil {
		t.Errorf(err.Error())
	}
	res, err := client.GetAndDelete(ctx, "tbb")
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("result is %s", res)
}
