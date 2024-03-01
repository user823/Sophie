package hash

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
		return newFnvHasher()
	}
}
