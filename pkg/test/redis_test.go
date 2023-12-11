package test

import (
	"context"
	"github.com/user823/Sophie/pkg/db/kv"
	"reflect"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestDisableRedis(t *testing.T) {
	type args struct {
		disable bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv.DisableRedis(tt.args.disable)
		})
	}
}

func TestConnected(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kv.Connected(); got != tt.want {
				t.Errorf("Connected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnectToRedis(t *testing.T) {
	type args struct {
		config *kv.RedisConfig
	}
	tests := []struct {
		name string
		args args
		want redis.UniversalClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kv.ConnectToRedis(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnectToRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeepConnection(t *testing.T) {
	type args struct {
		ctx    context.Context
		config *kv.RedisConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv.KeepConnection(tt.args.ctx, tt.args.config)
		})
	}
}

func TestDisconnect(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := kv.Disconnect(); (err != nil) != tt.wantErr {
				t.Errorf("Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRedisClient(t *testing.T) {
	tests := []struct {
		name string
		want kv.KeyValueStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kv.NewRedisClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRedisClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_Connect(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.Connect(tt.args.ctx, tt.args.config)
		})
	}
}

func Test_redisClient_Connected(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.Connected(); got != tt.want {
				t.Errorf("redisClient.Connected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_Disconnect(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.Disconnect(); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisClient_SetKeyPrefix(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		prefix string
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.SetKeyPrefix(tt.args.prefix)
		})
	}
}

func Test_redisClient_SetHashKey(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	type args struct {
		ok bool
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.SetHashKey(tt.args.ok)
		})
	}
}

func Test_redisClient_GetKey(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetKey(tt.args.ctx, tt.args.keyname)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetMultiKey(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetMultiKey(tt.args.ctx, tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetMultiKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetMultiKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetKeyTTL(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetKeyTTL(tt.args.ctx, tt.args.keyname)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetKeyTTL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.GetKeyTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetRawKey(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetRawKey(tt.args.ctx, tt.args.keyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetRawKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.GetRawKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetExp(t *testing.T) {
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
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetExp(tt.args.ctx, tt.args.keyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetExp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.GetExp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetKeys(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.GetKeys(tt.args.ctx, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_DeleteKey(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteKey(tt.args.ctx, tt.args.key); got != tt.want {
				t.Errorf("redisClient.DeleteKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_DeleteAllKeys(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteAllKeys(tt.args.ctx); got != tt.want {
				t.Errorf("redisClient.DeleteAllKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_DeleteRawKey(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteRawKey(tt.args.ctx, tt.args.key); got != tt.want {
				t.Errorf("redisClient.DeleteRawKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetKeysAndValuesWithFilter(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.GetKeysAndValuesWithFilter(tt.args.ctx, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetKeysAndValuesWithFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetKeysAndValues(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.GetKeysAndValues(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetKeysAndValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_DeleteKeys(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteKeys(tt.args.ctx, tt.args.keys); got != tt.want {
				t.Errorf("redisClient.DeleteKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_Decrement(t *testing.T) {
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.Decrement(tt.args.ctx, tt.args.key)
		})
	}
}

func Test_redisClient_IncrememntWithExpire(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.IncrememntWithExpire(tt.args.ctx, tt.args.key, tt.args.expire); got != tt.want {
				t.Errorf("redisClient.IncrememntWithExpire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_SetRollingWindow(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, got1 := r.SetRollingWindow(tt.args.ctx, tt.args.key, tt.args.per, tt.args.val, tt.args.pipeline)
			if got != tt.want {
				t.Errorf("redisClient.SetRollingWindow() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("redisClient.SetRollingWindow() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_redisClient_GetRollingWindow(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, got1 := r.GetRollingWindow(tt.args.ctx, tt.args.key, tt.args.per, tt.args.pipeline)
			if got != tt.want {
				t.Errorf("redisClient.GetRollingWindow() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("redisClient.GetRollingWindow() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_redisClient_GetKeyPrefix(t *testing.T) {
	type fields struct {
		KeyPrefix string
		HashKey   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.GetKeyPrefix(); got != tt.want {
				t.Errorf("redisClient.GetKeyPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetSet(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetSet(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_AddToSet(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.AddToSet(tt.args.ctx, tt.args.key, tt.args.value)
		})
	}
}

func Test_redisClient_GetAndDeleteSet(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.GetAndDeleteSet(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetAndDeleteSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_RemoveFromSet(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.RemoveFromSet(tt.args.ctx, tt.args.key, tt.args.value)
		})
	}
}

func Test_redisClient_DeleteScanMatch(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if got := r.DeleteScanMatch(tt.args.ctx, tt.args.filter); got != tt.want {
				t.Errorf("redisClient.DeleteScanMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_AddToSortedSet(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.AddToSortedSet(tt.args.ctx, tt.args.key, tt.args.value, tt.args.score)
		})
	}
}

func Test_redisClient_GetSortedSetRange(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, got1, err := r.GetSortedSetRange(tt.args.ctx, tt.args.key, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetSortedSetRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetSortedSetRange() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("redisClient.GetSortedSetRange() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_redisClient_RemoveSortedSetRange(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.RemoveSortedSetRange(tt.args.ctx, tt.args.key, tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.RemoveSortedSetRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisClient_GetListRange(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.GetListRange(tt.args.ctx, tt.args.key, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetListRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetListRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_RemoveFromList(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.RemoveFromList(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.RemoveFromList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisClient_AppendToSet(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			r.AppendToSet(tt.args.ctx, tt.args.key, tt.args.value)
		})
	}
}

func Test_redisClient_Exists(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			got, err := r.Exists(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_SetExp(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.SetExp(tt.args.ctx, tt.args.key, tt.args.expire); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.SetExp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisClient_SetKey(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.SetKey(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expire); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.SetKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisClient_SetRawKey(t *testing.T) {
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
			r := &kv.redisClient{
				KeyPrefix: tt.fields.KeyPrefix,
				HashKey:   tt.fields.HashKey,
			}
			if err := r.SetRawKey(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expire); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.SetRawKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
