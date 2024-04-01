package hash

import (
	"encoding/hex"
	"hash"
	"hash/crc32"
)

type crc32Hasher struct {
	crc hash.Hash
}

func newCrc32Hasher() Hasher {
	tb := crc32.MakeTable(crc32.Koopman)
	crc := crc32.New(tb)
	return &crc32Hasher{
		crc: crc,
	}
}

func (c *crc32Hasher) HashKey(key []byte) string {
	c.crc.Reset()
	c.crc.Write([]byte(key))
	return hex.EncodeToString(c.crc.Sum(nil))
}
