package models

import (
	"fmt"
	"net"
)

const NodesPrefix = "sophie-schedule-nodes"

const SyncTimeout = 3

// 需要注册的节点信息
type Node struct {
	// uuid
	Id        string `json:"id"`
	Network   string `json:"network"`
	IPAddress string `json:"ip_address"`
	// 负载权重信息（比如从节点中可定义为当前实例job数)
	LoadWeight int `json:"load_weight"`
	// 节点类型
	Mode string `json:"mode"`
	// 节点状态
	Status string `json:"status"`
	// 描述信息
	Description string `json:"description,omitempty"`
	// 标签信息
	Tags map[string]string `json:"tags,omitempty"`
}

// 节点池类型, 包含了所有注册的节点信息，便于统一管理
// 它只是etcd node视图，不能操作etcd中的数据
type NodePool interface {
	// 返回node数目
	Size() int
	// 获取所有在线节点
	OnlineNodes() []Node
	// 获取所有节点
	AllNodes() []Node
	// 根据id 获取node
	GetNode(id string) (Node, error)
}

// job 调度信息
type JobMap interface {
	// 查询job 是否存在；如果存在，那么返回在哪个节点上
	Exists(int64) (string, bool)
	// 查询全局job 试图
	Jobs() map[int64]string
	// 删除job调度信息
	Delete(int64) error
	// 创建job调度信息
	Create(int64, string) error
	// 查询node上运行的任务
	JobsOnNode(nodeid string) []int64
}

// 节点实例类型
type Instance interface {
	CacheKey() string
	Address() net.Addr
	Weight() int
	Tag(string) (value string, exist bool)
}

func (n Node) Address() net.Addr {
	switch n.Network {
	case "tcp", "tcp4", "tcp6":
		nt, _ := net.ResolveTCPAddr(n.Network, n.IPAddress)
		return nt
	case "udp", "udp4", "udp6":
		nt, _ := net.ResolveUDPAddr(n.Network, n.IPAddress)
		return nt
	case "ip", "ip4", "ip6":
		nt, _ := net.ResolveIPAddr(n.Network, n.IPAddress)
		return nt
	}
	return nil
}

func (n Node) Weight() int {
	return n.LoadWeight
}

func (n Node) Tag(key string) (value string, exist bool) {
	value, exist = n.Tags[key]
	return
}

func (n Node) CacheKey() string {
	return n.Id
}

func ServiceKey(cacheKey string) string {
	return fmt.Sprintf("%s:%s", NodesPrefix, cacheKey)
}

//sophie-schedule-nodes:449098f7-699f-4794-a8ee-69495414557e
//{"id":"449098f7-699f-4794-a8ee-69495414557e","network":"tcp","ip_address":"127.0.0.1:8091","load_weight":1000,"mode":"worker","status":"online","description":"449098f7-699f-4794-a8ee-69495414557e node worker, created at 2024-03-27T20:16:54+08:00"}
