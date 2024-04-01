package aggregation

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// redis 前缀
	RecordPrefix = "sophie_aggregation-"
	// 中间管道key
	Recordkey = "sophie_records"
	// 输出端key
	RecordAggregation                 = "sophie_record_aggregation"
	MessageTag                        = "aggregation"
	recoredsBufferForcedFlushInterval = 1 * time.Second
	max_size                          = 5 * 1024 * 1024
)

// 日志聚合类型
type LogRecord struct {
	Level      string         `json:"level,omitempty" mapstructure:"level,omitempty"`
	Time       string         `json:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`
	LoggerName string         `json:"logger,omitempty" mapstructure:"logger,omitempty"`
	Message    string         `json:"message,omitempty" mapstructure:"message,omitempty"`
	Caller     string         `json:"caller,omitempty" mapstructure:"caller,omitempty"`
	Stack      string         `json:"stack,omitempty" mapstructure:"stack,omitempty"`
	Additional map[string]any `json:"additional,omitempty" mapstructure:"additional,remain"`
}

// 日志发送端
type RecordProducer interface {
	Connect() bool
	AppendToSet(ctx context.Context, msg []string)
	Stop() error
}

// 日志接受端
type RecordConsumer interface {
	Connect() bool
	GetAndDeleteSet(ctx context.Context) []LogRecord
	Stop() error
}

var (
	analytics *Analytics
)

type Analytics struct {
	producer         RecordProducer
	poolSize         int
	workerBufferSize uint64
	recordBufferSize uint64
	// 单位毫秒
	recordsBufferFlushInterval uint64
	poolWg                     sync.WaitGroup
	flush                      atomic.Bool
	recordDetail               bool
	recordCh                   chan string
	shouldStop                 uint32
}

func init() {
	analytics = &Analytics{}
}

func NewAnalytics(options *AnalyticsOptions, p RecordProducer) {
	poolsize := options.PoolSize
	recordsBufferSize := options.RecordsBufferSize
	workerBufferSize := recordsBufferSize / uint64(poolsize)
	analytics.recordCh = make(chan string, recordsBufferSize)
	analytics.producer = p
	analytics.poolSize = poolsize
	analytics.workerBufferSize = workerBufferSize
	analytics.recordDetail = options.EnableDetailedRecording
	analytics.recordBufferSize = recordsBufferSize
	analytics.recordsBufferFlushInterval = options.FlushInterval
}

func GetAnalytics() *Analytics {
	return analytics
}

func (r *Analytics) Start() {
	// 等待连接建立完成
	for !r.producer.Connect() {
		time.Sleep(1 * time.Second)
	}

	atomic.SwapUint32(&r.shouldStop, 0)

	time.Sleep(3 * time.Second)
	for i := 0; i < r.poolSize; i++ {
		r.poolWg.Add(1)
		go r.recordWorker()
	}
}

func (r *Analytics) Stop() {
	atomic.SwapUint32(&r.shouldStop, 1)

	// 等待组件优雅关停
	if err := r.producer.Stop(); err != nil {
		// 不处理
	}
	close(r.recordCh)
	// 等待剩余的任务输出完毕
	r.poolWg.Wait()
}

func (r *Analytics) Write(data []byte) (int, error) {
	if r.recordCh == nil {
		return len(data), nil
	}

	// 丢弃过大数据
	if !r.recordDetail && len(data) > max_size {
		return len(data), nil
	}

	// 将拷贝后的数据写入通道
	r.recordCh <- string(data)
	return len(data), nil
}

func (r *Analytics) Sync() error {
	r.flush.Store(true)
	return nil
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
		case record, ok := <-r.recordCh:
			// 通道关闭
			if !ok {
				r.producer.AppendToSet(context.Background(), recordsBuffer)
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
			r.producer.AppendToSet(context.Background(), recordsBuffer)
			recordsBuffer = recordsBuffer[:0]
			lastSent = time.Now()
		}
	}
}
