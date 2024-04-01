package hash

import (
	"encoding/hex"
	"hash"
	"hash/maphash"
)

type mapHasher struct {
	h hash.Hash
}

func newMapHasher() Hasher {
	return &mapHasher{h: &maphash.Hash{}}
}

func (m *mapHasher) HashKey(key []byte) string {
	m.h.Reset()
	m.h.Write(key)
	return hex.EncodeToString(m.h.Sum(nil))
}
