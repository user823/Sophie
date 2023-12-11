package hash

import "github.com/user823/Sophie/pkg/log"

const DefaultHashAlgorithm = "fnv"

type Hasher interface {
	HashKey(string) string
}

func NewHasher(strategy string) Hasher {
	switch strategy {
	case "crc":
		return newCrc32Hasher()
	case "adler":
		return newAdler32Hasher()
	case "fnv":
		return newFnvHasher()
	case "maphash":
		return newMapHasher()
	default:
		log.Debugln("Default choose fnv hash algo.")
		return newFnvHasher()
	}
}
