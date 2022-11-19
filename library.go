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

type Register struct {
	Id   string
	Data interface{}
}

// Prefix and suffix handling:
//  0. prefix_suffix ...string is [prefix, suffix, ignored, ignored...]string
//     in other words, prefix==prefix_suffix[0], suffix==prefix_suffix[1]
//  1. If NO PREFIX is specified, the path will be ./person/
//  2. If PREFIX=data, the path will be ./data/person/
//  3. If SUFFIX=people, the path will be ./data/people/
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
	gx.Trace("JDocDB INSERT: ", id, ": ", jsonPath)
}

// Selects one registry from a table using its ID, prefix is a set of dir/subdirectories
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
	gx.Trace("JDocDB SELECT: ", id, ": ", jsonPath)
	return doc
}

// Select all IDs of a table, prefix is a set of dir/subdirectories
func SelectIds[T interface{}](doc T, prefix ...string) []string {
	idList := []string{}
	table := buildPath(GetType(doc), prefix...)
	fileList, err := ioutil.ReadDir(table)
	gx.Error(err)
	for _, f := range fileList {
		if strings.HasSuffix(f.Name(), ".json") {
			gx.Trace("JDocDB SELECT_IDS: found: ", f.Name())
			idList = append(idList, strings.TrimSuffix(f.Name(), ".json"))
		}
	}
	gx.Trace("JDocDB SELECT_IDS: ", idList, ": ", table)
	return idList
}

// Selects all rows from a table, prefix is a set of dir/subdirectories
func SelectAll[T interface{}](doc T, prefix ...string) map[string]T {
	docs := map[string]T{}
	for _, id := range SelectIds(doc, prefix...) {
		docs[id] = Select(id, doc, prefix...)
	}
	gx.Trace("JDocDB SELECT_ALL: ", docs)
	return docs
}

// Selects all rows that meet some conditions, prefix is a set of dir/subdirectories
func SelectWhere[T interface{}](doc T, cond func(T) bool, prefix ...string) map[string]T {
	docs := map[string]T{}
	for _, id := range SelectIds(doc, prefix...) {
		candidate := Select(id, doc, prefix...)
		if cond(candidate) {
			docs[id] = candidate
		}
	}
	gx.Trace("JDocDB SELECT_WHERE: ", docs)
	return docs
}

// Selects all rows that meet some conditions, prefix is a set of dir/subdirectories
func SelectIdWhere[T interface{}](doc T, cond func(T) bool, prefix ...string) []string {
	idList := []string{}
	for _, id := range SelectIds(doc, prefix...) {
		candidate := Select(id, doc, prefix...)
		if cond(candidate) {
			idList = append(idList, id)
		}
	}
	gx.Trace("JDocDB SELECT_ID_WHERE: ", idList)
	return idList
}

// Returns a lowercase string with the type
func GetType(doc interface{}) string {
	return strings.ToLower(strings.SplitN(fmt.Sprintf("%T", doc), ".", 2)[1])
}
