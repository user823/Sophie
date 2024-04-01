package loadbalance

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/hash"
	"sort"
)

type ConsistentHashPicker struct {
	instances []models.Instance
	keys      []string
	mp        map[string]int
	// 虚拟节点数目
	replica int
	hasher  hash.Hasher
}

func newConsistentHashPicker(instances []models.Instance, hasher hash.Hasher, replica int) Picker {
	if replica <= 0 {
		replica = 1
	}
	if hasher == nil {
		hasher = hash.NewHasher(hash.DefaultHashAlgorithm)
	}
	c := new(ConsistentHashPicker)
	c.mp = map[string]int{}
	c.replica = replica
	c.hasher = hasher
	c.instances = instances

	for idx, instance := range instances {
		key := instance.CacheKey()
		for i := 0; i < c.replica; i++ {
			// 对于多副本的情况，使用 key_index 的格式构造cachekey
			cacheKey := fmt.Sprintf("%s_%d", key, i)
			hashCode := c.hasher.HashKey(utils.S2b(cacheKey))
			if _, ok := c.mp[hashCode]; !ok {
				c.keys = append(c.keys, hashCode)
				c.mp[hashCode] = idx
			}
		}
	}
	sort.Strings(c.keys)
	return c
}

func (c *ConsistentHashPicker) Next(ctx context.Context, req any) models.Instance {
	if len(c.keys) == 0 {
		return nil
	}

	// 只支持string类型
	cacheKey, ok := req.(string)
	if !ok {
		return nil
	}

	hashCode := c.hasher.HashKey(utils.S2b(cacheKey))
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hashCode
	})

	if idx == len(c.keys) {
		idx = 0
	}
	selectedNode := c.mp[c.keys[idx]]
	return c.instances[selectedNode]
}
