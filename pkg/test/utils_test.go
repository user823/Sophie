package test

import (
	"fmt"
	"github.com/user823/Sophie/pkg/utils/hash"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"testing"
)

func TestHash(t *testing.T) {
	hasher := hash.NewHasher(hash.DefaultHashAlgorithm)
	fmt.Println(hasher.HashKey("123"))
	fmt.Println(hasher.HashKey("123"))
	fmt.Println(hasher.HashKey("12"))
}

func TestSimpleMatch(t *testing.T) {
	args := [][]string{
		{"*12345*", "123"},
		{"*12345*", "12345"},
		{"*12345", "12345"},
		{"12345*", "12345"},
		{"123456", "12345"},
		{"*1123*", "12345"},
	}
	wants := []bool{
		false,
		true,
		true,
		true,
		false,
		false,
	}
	for i, arg := range args {
		if res := strutil.SimpleMatch(arg[0], arg[1]); wants[i] != res {
			t.Errorf("want %v but got %v, args: %s %s\n", wants[i], res, arg[0], arg[1])
		}
	}
}
