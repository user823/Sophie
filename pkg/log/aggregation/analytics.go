package aggregation

import (
	"context"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"go.uber.org/zap/zapcore"
	"sync"
	"sync/atomic"
	"time"
)

const (
	RecordkeyName                     = "record"
	recoredsBufferForcedFlushInterval = 1 * time.Second
	max_size                          = 5 * 1024 * 1024
)

// 后续用于promp 聚合到elasticsearch中
type LogRecord struct {
	Level      zapcore.Level  `json:"level,omitempty" mapstructure:"level,omitempty"`
	Time       time.Time      `json:"time,omitempty" mapstructure:"timestamp,omitempty"`
	LoggerName string         `json:"logger,omitempty" mapstructure:"logger,omitempty"`
	Message    string         `json:"message,omitempty" mapstructure:"message,omitempty"`
	Caller     string         `json:"caller,omitempty" mapstructure:"caller,omitempty"`
	Stack      string         `json:"stack,omitempty" mapstructure:"stack,omitempty"`
	Additional map[string]int `json:"additional,omitempty" mapstructure:"additional,remain"`
}

var (
	analytics *Analytics
	once      sync.Once
)

type Analytics struct {
	store            kv.RedisStore
	poolSize         int
	workerBufferSize uint64
	recordBufferSize uint64
	// 单位毫秒
	recordsBufferFlushInterval uint64
	poolWg                     sync.WaitGroup
	rchManager                 RecordChManager
	flush                      atomic.Bool
	recordDetail               bool
}

type RecordChManager interface {
	Start(bufSize uint64)
	Stop()
	ShouldStop() bool
	GetChannel() chan string
}

func NewAnalytics(options *AnalyticsOptions, store kv.RedisStore, rchManager RecordChManager) {
	once.Do(func() {
		poolsize := options.PoolSize
		recordsBufferSize := options.RecordsBufferSize
		workerBufferSize := recordsBufferSize / uint64(poolsize)
		log.Debug("Analytics pool worker buffer size", "workerBufferSize", workerBufferSize)
		store.SetKeyPrefix(kv.LogKeyPrefix)
		store.SetHashKey(true)
		analytics = &Analytics{
			store:                      store,
			poolSize:                   poolsize,
			workerBufferSize:           workerBufferSize,
			recordsBufferFlushInterval: options.FlushInterval,
			rchManager:                 rchManager,
			recordBufferSize:           recordsBufferSize,
			recordDetail:               options.EnableDetailedRecording,
		}
	})
}

func GetAnalytics() *Analytics {
	return analytics
}

func (r *Analytics) Start() {
	// 等待连接建立完成
	for !r.store.Connected() {
		log.Info("waiting for redis connect establish")
		time.Sleep(1 * time.Second)
	}
	r.rchManager.Start(r.recordBufferSize)
	for i := 0; i < r.poolSize; i++ {
		r.poolWg.Add(1)
		go r.recordWorker()
	}
}

func (r *Analytics) Stop() {
	r.rchManager.Stop()
	r.poolWg.Wait()
}

func (r *Analytics) RecordHit(record string) {
	if r.rchManager.ShouldStop() {
		return
	}

	// 超过记录最大长度
	if !r.recordDetail && len(record) > max_size {
		return
	}
	r.rchManager.GetChannel() <- record
}

func (r *Analytics) Flush() {
	r.flush.Store(true)
}

func (r *Analytics) recordWorker() {
	defer r.poolWg.Done()
	recordsBuffer := make([]string, 0, r.workerBufferSize)
	var readyToSend bool
	lastSent := time.Now()
	for {
		readyToSend = false
		select {
		case record, ok := <-r.rchManager.GetChannel():
			// 通道关闭
			if !ok {
				r.store.AppendToSetPipelined(context.Background(), RecordkeyName, recordsBuffer)
				return
			}

			recordsBuffer = append(recordsBuffer, record)
			readyToSend = uint64(len(recordsBuffer)) == r.workerBufferSize
		case <-time.After(time.Duration(r.recordsBufferFlushInterval) * time.Millisecond):
			readyToSend = true
		default:
			// flush
			if r.flush.CompareAndSwap(true, false) {
				readyToSend = true
			}
		}
		if len(recordsBuffer) > 0 && (readyToSend || time.Since(lastSent) >= recoredsBufferForcedFlushInterval) {
			r.store.AppendToSetPipelined(context.Background(), RecordkeyName, recordsBuffer)
			recordsBuffer = recordsBuffer[:0]
			lastSent = time.Now()
		}
	}
}
