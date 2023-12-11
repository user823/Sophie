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

func (m *mapHasher) HashKey(keyname string) string {
	m.h.Reset()
	m.h.Write([]byte(keyname))
	return hex.EncodeToString(m.h.Sum(nil))
}
