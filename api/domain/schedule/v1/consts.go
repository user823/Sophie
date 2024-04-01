package v1

const (
	ServiceName = "Sophie Schedule"
)

// 节点状态
const (
	NODE_STATUS_ON = "online"
	// worker节点可以主动更新自己状态为离线
	NODE_STATUS_OFF = "offline"
	// worker节点可以主动更新自己状态为挂起状态，此时manager不要把它的任务调度到其他节点上
	// 此时也不应该接受调度
	NODE_STATUS_PAUSE = "pause"
)

// 任务调度状态常量
const (
	NORMAL = "0"
	PAUSE  = "1"
)

// 节点类型
const (
	MANAGER = "manager"
	WORKER  = "worker"
)

// etcd 分布式锁
const ()
