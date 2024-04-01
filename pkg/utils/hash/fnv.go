package hash

import (
	"encoding/hex"
	"hash"
	"hash/fnv"
)

type fnvHasher struct {
	f hash.Hash
}

func newFnvHasher() Hasher {
	return &fnvHasher{
		f: fnv.New32(),
	}
}

func (f *fnvHasher) HashKey(key []byte) string {
	f.f.Reset()
	f.f.Write([]byte(key))
	return hex.EncodeToString(f.f.Sum(nil))
}
