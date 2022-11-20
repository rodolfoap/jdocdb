package jdocdb

import (
	"encoding/json"
	"fmt"
	"github.com/rodolfoap/gx"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Structure used for unmarshaling each register. The Id value is generated and
// kept while recording, because it duplicates the file name.
type Register struct {
	Id   string
	Data interface{}
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

// Inserts one registry into a table using its ID, prefix is a set of dir/subdirectories
func Insert[T interface{}](id string, doc T, prefix ...string) {
	reg := Register{Id: id, Data: doc}
	table := buildPath(GetType(doc), prefix...)
	os.MkdirAll(table, 0755)
	jsonPath := filepath.Join(table, id+".json")
	jsonBytes, err := json.MarshalIndent(reg, "", "\t")
	gx.Error(err)
	jsonBytes = append(jsonBytes, byte('\n'))
	err = ioutil.WriteFile(jsonPath, jsonBytes, 0644)
	gx.Fatal(err)
	gx.Tracef("JDocDB INSERT: %v", jsonPath)
}

// Selects one registry from a table using its ID, prefix is a set of dir/subdirectories.
// This just reads a single file, unmarshals it and returns a (generic) structure.
func Select[T interface{}](id string, doc T, prefix ...string) T {
	reg := Register{Id: id, Data: &doc}
	table := buildPath(GetType(doc), prefix...)
	jsonPath := filepath.Join(table, id+".json")
	jsonBytes, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		gx.Tracef("Data file %v not found.", jsonPath)
		return doc
	}
	err = json.Unmarshal(jsonBytes, &reg)
	gx.Error(err)
	gx.Tracef("JDocDB SELECT: %v", jsonPath)
	return doc
}

// Select all IDs of a table, prefix is a set of dir/subdirectories.
// Internally, just finding, not unmarshaling, .json files in a directory,
// since IDs are just filename prefixes.
func SelectIds[T interface{}](doc T, prefix ...string) []string {
	idList := []string{}
	table := buildPath(GetType(doc), prefix...)
	fileList, err := ioutil.ReadDir(table)
	gx.Error(err)
	for _, f := range fileList {
		if strings.HasSuffix(f.Name(), ".json") {
			idList = append(idList, strings.TrimSuffix(f.Name(), ".json"))
		}
	}
	gx.Trace("JDocDB SELECT_IDS: ", table, idList)
	return idList
}

// Selects all rows from a table, prefix is a set of dir/subdirectories.
// Combines SelectIds(), which provides the list of files, and Select(), which
// unmarshals each file into the structure provided as a generic type.
func SelectAll[T interface{}](doc T, prefix ...string) map[string]T {
	docs := map[string]T{}
	for _, id := range SelectIds(doc, prefix...) {
		docs[id] = Select(id, doc, prefix...)
	}
	gx.Tracef("JDocDB SELECT_ALL: %v/%v", strings.Join(prefix, "/"), Keys(docs))
	return docs
}

// Selects all rows that meet some conditions, prefix is a set of dir/subdirectories.
// This applies the condition over each unmarshaled struct, providing the matching result in a map.
func SelectWhere[T interface{}](doc T, cond func(T) bool, prefix ...string) map[string]T {
	docs := map[string]T{}
	for key, val := range SelectAll(doc, prefix...) {
		if cond(val) {
			docs[key] = val
		}
	}
	gx.Trace("JDocDB SELECT_WHERE: ", docs)
	return docs
}

// Selects all rows that meet the conditions provided in the passing function, prefix
// is a set of dir/subdirectories.
func SelectIdWhere[T interface{}](doc T, cond func(T) bool, prefix ...string) []string {
	keys := Keys(SelectWhere(doc, cond, prefix...))
	gx.Trace("JDocDB SELECT_ID_WHERE: ", keys)
	return keys
}

// Selects all rows that meet some conditions, prefix is a set of dir/subdirectories
func SelectWhereGroup[T interface{}, A interface{}](doc T, cond func(T) bool, aggregator *A, aggregate func(T), prefix ...string) map[string]T {
	docs := SelectWhere(doc, cond, prefix...)
	for _, val := range docs {
		aggregate(val)
	}
	gx.Trace("JDocDB SELECT_ID_WHERE_GROUP: ", docs, *aggregator)
	return docs
}

// Returns a lowercase string with the type name for path finding.
func GetType(doc interface{}) string {
	return strings.ToLower(strings.SplitN(fmt.Sprintf("%T", doc), ".", 2)[1])
}

// Keys returns the keys of the map m. Keys will be in an indeterminate order. Taken
// from https://github.com/golang/exp/blob/master/maps/maps.go which is subject to change
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
