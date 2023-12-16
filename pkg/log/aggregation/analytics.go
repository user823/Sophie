package aggregation

import (
	"context"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"go.uber.org/zap/zapcore"
	"sync"
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
	workerBufferSize uint64
	recordBufferSize uint64
	// 单位毫秒
	recordsBufferFlushInterval uint64
	poolWg                     sync.WaitGroup
	rchManager                 *log.RecordChManager
}

func NewAnalytics(options *AnalyticsOptions, store kv.RedisStore, rchManager *log.RecordChManager) {
	once.Do(func() {
		poolsize := options.PoolSize
		recordsBufferSize := options.RecordsBufferSize
		workerBufferSize := recordsBufferSize / uint64(poolsize)
		log.Debug("Analytics pool worker buffer size", "workerBufferSize", workerBufferSize)
		store.SetKeyPrefix(LogKeyPrefix)
		store.SetHashKey(true)
		analytics = &Analytics{
			store:                      store,
			poolSize:                   poolsize,
			workerBufferSize:           workerBufferSize,
			recordsBufferFlushInterval: options.FlushInterval,
			rchManager:                 rchManager,
			recordBufferSize:           recordsBufferSize,
		}
	})
}

func GetAnalytics() *Analytics {
	return analytics
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
	r.rchManager.RecordCh <- record
}

func (r *Analytics) recordWorker() {
	defer r.poolWg.Done()
	recordsBuffer := make([]string, 0, r.workerBufferSize)
	var readyToSend bool
	lastSent := time.Now()
	for {
		readyToSend = false
		select {
		case record, ok := <-r.rchManager.RecordCh:
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
