package loadbalance

import (
	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils/hash"
	"github.com/user823/Sophie/pkg/utils/strutil"
)

type LoadBalanceOptions struct {
	Strategy string `json:"strategy" mapstructure:"strategy"`
	// 一致性hash算法的hash 函数
	Hash string `json:"hash" mapstructure:"hash"`
	// 使用一致性hash算法时 每个节点的虚拟节点数目
	HashNodes int `json:"hash_nodes" mapstructure:"hash_nodes"`
}

func NewLoadBalanceOptions() *LoadBalanceOptions {
	return &LoadBalanceOptions{
		Strategy:  WeightedRoundRobin,
		Hash:      hash.DefaultHashAlgorithm,
		HashNodes: 3,
	}
}

func (o *LoadBalanceOptions) Validate() error {
	if !strutil.ContainsAny(o.Strategy, WeightedRandom, WeightedRoundRobin, ConsistentHash) {
		return errors.New("loadbalancer strategy is invalid, only support 'weight_round_robin', 'weight_random' and 'consistent_hash'.")
	}
	return nil
}

func (o *LoadBalanceOptions) Complete() error {
	if o.HashNodes <= 0 {
		o.HashNodes = 3
	}
	return nil
}

func (o *LoadBalanceOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.Strategy, "loadbalance.strategy", o.Strategy, ""+
		"Choose loadbalance strategy, only supported 'weight_round_robin'(default), 'weight_random' and 'consistent_hash'")
	fs.StringVar(&o.Hash, "loadbalance.hash", o.Hash, ""+
		"Choose hash function to hash cachekey, supported 'crc', 'adler', 'fnv'(default), 'maphash'")
	fs.IntVar(&o.HashNodes, "loadbalance.hash_nodes", o.HashNodes, ""+
		"Set virtual nodes num for each really node, default 3")
}

func (o *LoadBalanceOptions) CreateLoadBalancer() LoadBalancer {
	switch o.Strategy {
	case WeightedRandom:
		return NewWeightedRandomBalancer()
	case WeightedRoundRobin:
		return NewWeightedRoundRoinBalancer()
	case ConsistentHash:
		hasher := hash.NewHasher(o.Strategy)
		return NewConsistentHashBalancer(hasher, o.HashNodes)
	}
	log.Errorf("loadbalancer options parameters is invalid")
	return nil
}
