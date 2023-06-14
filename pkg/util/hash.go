package util

import (
	"hash"
	"hash/fnv"

	"github.com/davecgh/go-spew/spew"
)

func Hash(object interface{}) uint64 {
	hash := fnv.New32a()
	deepHashObject(hash, object)
	return uint64(hash.Sum32())
}

// DeepHashObject writes specified object to hash using the spew library
// which follows pointers and prints actual values of the nested objects
// ensuring the hash does not change when a pointer changes.
func deepHashObject(hasher hash.Hash, objectToWrite interface{}) {
	hasher.Reset()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}
	printer.Fprintf(hasher, "%#v", objectToWrite)
}
