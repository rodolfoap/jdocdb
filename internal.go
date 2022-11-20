// A Golang JSON document-file based database allowing complex SELECT queries.
package jdocdb

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Returns a lowercase string with the type name for path finding.
func getType(doc interface{}) string {
	return strings.ToLower(strings.SplitN(fmt.Sprintf("%T", doc), ".", 2)[1])
}

// Keys returns the keys of the map m. Keys will be in an indeterminate order. Taken
// from https://github.com/golang/exp/blob/master/maps/maps.go which is subject to change
func keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// Prefix and suffix handling:
//
//	a) prefix_suffix ...string is [prefix, suffix, ignored, ignored...]string
//	   in other words, prefix==prefix_suffix[0], suffix==prefix_suffix[1].
//	b) If NO PREFIX is specified, the path will be ./person/;
//	c) If PREFIX=data, the path will be ./data/person/;
//	d) If SUFFIX=people, the path will be ./data/people/.
func buildPath(baseName string, prefix_suffix ...string) string {
	prefix, suffix := "", baseName
	if len(prefix_suffix) > 0 {
		prefix = prefix_suffix[0]
	}
	if len(prefix_suffix) > 1 {
		suffix = prefix_suffix[1]
	}
	buildPath := filepath.Clean(filepath.Join(prefix, suffix)) + "/"
	return buildPath
}
