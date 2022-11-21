// A Golang JSON document-file based database allowing complex SELECT queries and multiple aggregations.
package jdocdb

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/rodolfoap/gx"
)

// Structure used for unmarshaling each register. The Id value is generated and
// kept while recording, because it duplicates the file name.
type Register struct {
	Id   string
	Data interface{}
}

// Inserts one registry into a table using its ID, prefix is a set of dir/subdirectories
func Insert[T interface{}](id string, doc T, prefix ...string) {
	reg := Register{Id: id, Data: doc}
	table := buildPath(getType(doc), prefix...)
	os.MkdirAll(table, 0755)
	jsonPath := filepath.Join(table, id+".json")
	jsonBytes, err := json.MarshalIndent(reg, "", "\t")
	gx.Error(err)
	jsonBytes = append(jsonBytes, byte('\n'))
	err = os.WriteFile(jsonPath, jsonBytes, 0644)
	gx.Fatal(err)
	gx.Tracef("JDocDB INSERT: %v", jsonPath)
}

// Selects one registry from a table using its ID, prefix is a set of dir/subdirectories.
// This just reads a single file, unmarshals it and returns a (generic) structure.
func Select[T interface{}](id string, doc T, prefix ...string) T {
	reg := Register{Id: id, Data: &doc}
	table := buildPath(getType(doc), prefix...)
	jsonPath := filepath.Join(table, id+".json")
	jsonBytes, err := os.ReadFile(jsonPath)
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
	table := buildPath(getType(doc), prefix...)
	fileList, err := os.ReadDir(table)
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
	gx.Tracef("JDocDB SELECT_ALL: %v/%v", strings.Join(prefix, "/"), keys(docs))
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
	keys := keys(SelectWhere(doc, cond, prefix...))
	gx.Trace("JDocDB SELECT_ID_WHERE: ", keys)
	return keys
}

// Selects all rows that meet some conditions, prefix is a set of dir/subdirectories.
// Parameters:
//
//	doc T: Any user-defined type struct, for example, {Person}
//	cond func(): A user-defined type taking the previous struct and yielding a boolean. Every entry is compared to consider is it is returned.
//	aggregator *A: A pointer reference to a user variable, which is available during the internal processing loop, to calculate aggregates. Could be a slice or a struct.
//	aggregate func(): A user-defined function that fills up the aggregator variable(s).
//	prefix...: Only used to read the prefix (table location, or ./) and the suffix (table directory or lowecase(type)).
func SelectWhereAggreg[T interface{}, A interface{}](doc T, cond func(T) bool, aggregator *A, aggregate func(string, T), prefix ...string) map[string]T {
	_docs := SelectWhere(doc, cond, prefix...)
	for _key, _val := range _docs {
		aggregate(_key, _val)
	}
	gx.Trace("JDocDB SELECT_WHERE_AGGREG: ", _docs, *aggregator)
	return _docs
}

// Equivalent to SelectWhereAggreg() except without WHERE clause.
func SelectAggreg[T interface{}, A interface{}](doc T, aggregator *A, aggregate func(string, T), prefix ...string) map[string]T {
	_docs := SelectAll(doc, prefix...)
	for _key, _val := range _docs {
		aggregate(_key, _val)
	}
	gx.Trace("JDocDB SELECT_AGGREG: ", _docs, *aggregator)
	return _docs
}

// Equivalent to SelectWhereAggreg() except that it returns just a count
func CountWhereAggreg[T interface{}, A interface{}](doc T, cond func(T) bool, aggregator *A, aggregate func(string, T), prefix ...string) int {
	_docs, _count := SelectWhere(doc, cond, prefix...), 0
	for key, val := range _docs {
		_count += 1
		aggregate(key, val)
	}
	gx.Trace("JDocDB COUNT_WHERE_AGGREG: ", _docs, *aggregator)
	return _count
}

// Equivalent to SelectWhereAggreg() except that it returns just a count
func CountWhere[T interface{}](doc T, cond func(T) bool, prefix ...string) int {
	docs := SelectWhere(doc, cond, prefix...)
	gx.Trace("JDocDB COUNT_WHERE: ", len(docs))
	return len(docs)
}

// Equivalent to CountWhereAggreg() except without WHERE conditionals
func CountAggreg[T interface{}, A interface{}](doc T, aggregator *A, aggregate func(string, T), prefix ...string) int {
	_docs, _count := SelectAll(doc, prefix...), 0
	for key, val := range _docs {
		_count += 1
		aggregate(key, val)
	}
	gx.Trace("JDocDB COUNT_AGGREG: ", _docs, *aggregator)
	return _count
}

// Simple count of all registers
func Count[T interface{}](doc T, prefix ...string) int {
	docs := SelectAll(doc, prefix...)
	gx.Trace("JDocDB COUNT: ", len(docs))
	return len(docs)
}

// Simple sum of all registers
func Sum[T interface{}](doc T, fieldName string, prefix ...string) int {
	_docs, _sum := SelectAll(doc, prefix...), 0
	var data []int
	for _, doc := range _docs {
		data = append(data, int(reflect.Indirect(reflect.ValueOf(doc)).FieldByName(fieldName).Int()))
	}
	gx.Trace("JDocDB SUM: ", _sum, data)
	// Getting an array would simplify a lot any future aggregate calculation.
	return sum(data)
}

// Simple sum of all registers fulfilling a WHERE condition
func SumWhere[T interface{}](doc T, fieldName string, cond func(T) bool, prefix ...string) int {
	_docs, _sum := SelectWhere(doc, cond, prefix...), 0
	var data []int
	for _, doc := range _docs {
		data = append(data, int(reflect.Indirect(reflect.ValueOf(doc)).FieldByName(fieldName).Int()))
	}
	gx.Trace("JDocDB SUM: ", _sum, data)
	// Getting an array would simplify a lot any future aggregate calculation.
	return sum(data)
}

// Deletes one registry from a table using its ID, prefix is a set of dir/subdirectories
func Delete[T interface{}](id string, doc T, prefix ...string) {
	table := buildPath(getType(doc), prefix...)
	jsonPath := filepath.Join(table, id+".json")
	err := os.Remove(jsonPath)
	gx.Fatal(err)
	gx.Tracef("JDocDB DELETE: %v", jsonPath)
}
