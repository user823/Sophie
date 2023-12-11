package log

import (
	"context"
	"github.com/user823/Sophie/pkg/db/kv"
	"go.uber.org/zap/zapcore"
	"sync"
	"sync/atomic"
	"time"
)

const (
	LogKeyPrefix                      = "sophie-log-"
	RecordkeyName                     = "record"
	recoredsBufferForcedFlushInterval = 1 * time.Second
)

type LogRecord struct {
	Level      zapcore.Level  `mapstructure:"level,omitempty"`
	Time       time.Time      `mapstructure:"timestamp,omitempty"`
	LoggerName string         `mapstructure:"logger,omitempty"`
	Message    string         `mapstructure:"message,omitempty"`
	Caller     string         `mapstructure:"caller,omitempty"`
	Stack      string         `mapstructure:"stacktrace,omitempty"`
	Additional map[string]int `mapstructure:",remain"`
}

var (
	analytics *Analytics
	once      sync.Once
)

type Analytics struct {
	store            kv.RedisStore
	poolSize         int
	recordsChan      chan string
	workerBufferSize uint64
	// 单位毫秒
	recordsBufferFlushInterval uint64
	shouldStop                 uint32
	poolWg                     sync.WaitGroup
}

func NewAnalytics(options *AnalyticsOptions, store kv.RedisStore) {
	once.Do(func() {
		poolsize := options.PoolSize
		recordsBufferSize := options.RecordsBufferSize
		workerBufferSize := recordsBufferSize / uint64(poolsize)
		Debug("Analytics pool worker buffer size", "workerBufferSize", workerBufferSize)
		recordsChan := make(chan string, recordsBufferSize)
		store.SetKeyPrefix(LogKeyPrefix)
		store.SetHashKey(true)
		analytics = &Analytics{
			store:                      store,
			poolSize:                   poolsize,
			recordsChan:                recordsChan,
			workerBufferSize:           workerBufferSize,
			recordsBufferFlushInterval: options.FlushInterval,
		}
	})
}

func GetAnalytics() *Analytics {
	return analytics
}

func Start() {
	Info("===> starting log aggregation system")
	analytics.Start()
}

func Stop() {
	Info("===> stopping log aggregation system")
	analytics.Stop()
}

func (r *Analytics) Start() {
	// 等待连接建立完成
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !r.store.Connected() {
		}
	}()
	wg.Wait()
	atomic.SwapUint32(&r.shouldStop, 0)
	for i := 0; i < r.poolSize; i++ {
		r.poolWg.Add(1)
		go r.recordWorker()
	}
}

func (r *Analytics) Stop() {
	atomic.SwapUint32(&r.shouldStop, 1)
	close(r.recordsChan)
	r.poolWg.Wait()
}

func (r *Analytics) RecordHit(record string) {
	if atomic.LoadUint32(&r.shouldStop) > 0 {
		return
	}
	r.recordsChan <- record
}

func (r *Analytics) recordWorker() {
	defer r.poolWg.Done()
	recordsBuffer := make([]string, 0, r.workerBufferSize)
	lastSent := time.Now()
	for {
		var readyToSend bool
		select {
		case record, ok := <-r.recordsChan:
			// 通道关闭
			if !ok {
				r.store.AppendToSetPipelined(context.Background(), RecordkeyName, recordsBuffer)
				return
			}

			recordsBuffer = append(recordsBuffer, record)
			readyToSend = uint64(len(recordsBuffer)) == r.workerBufferSize
		case <-time.After(time.Duration(r.recordsBufferFlushInterval) * time.Millisecond):
			readyToSend = true

		}

		if len(recordsBuffer) > 0 && (readyToSend || time.Since(lastSent) >= recoredsBufferForcedFlushInterval) {
			r.store.AppendToSetPipelined(context.Background(), RecordkeyName, recordsBuffer)
			recordsBuffer = recordsBuffer[:0]
			lastSent = time.Now()
		}
	}
}
