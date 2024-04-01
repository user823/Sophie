package hash

import (
	"encoding/hex"
	"hash"
)
import "hash/adler32"

type adler32Hasher struct {
	adler hash.Hash
}

func newAdler32Hasher() Hasher {
	return &adler32Hasher{
		adler: adler32.New(),
	}
}

func (a *adler32Hasher) HashKey(key []byte) string {
	a.adler.Reset()
	b := a.adler.Sum(key)
	return hex.EncodeToString(b)
}
