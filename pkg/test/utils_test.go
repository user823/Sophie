package test

import (
	"fmt"
	"github.com/user823/Sophie/pkg/utils/hash"
	"testing"
)

func TestHash(t *testing.T) {
	hasher := hash.NewHasher(hash.DefaultHashAlgorithm)
	fmt.Println(hasher.HashKey("123"))
	fmt.Println(hasher.HashKey("123"))
	fmt.Println(hasher.HashKey("12"))
}
