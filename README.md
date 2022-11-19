# jdocdb: A JSON File-documents Database
[![Go Reference](https://pkg.go.dev/badge/github.com/rodolfoap/jdocdb.svg)](https://pkg.go.dev/github.com/rodolfoap/jdocdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/rodolfoap/jdocdb)](https://goreportcard.com/report/github.com/rodolfoap/jdocdb)
[![Coverage Status](https://coveralls.io/repos/github/rodolfoap/jdocdb/badge.svg?branch=main)](https://coveralls.io/github/rodolfoap/jdocdb?branch=main)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
[![GPLv3 license](https://img.shields.io/badge/License-GPLv3-blue.svg)](http://perso.crans.org/besson/LICENSE.html)

A minimalist file-based JSON documents database with the capability of complex SELECT WHERE operations over a single table.

* Tables are subdirectories, e.g. `./clients/`;
* Registries are files, e.g. `./clients/a929782.json`;
* Filenames are registry IDs, e.g. `./clients/a929782.json` has `ID==a929782`;
* SQL SELECT equivalents are:
	* `Select(id, struct, tableLocation)`: "SELECT * FROM TABLE WHERE ID=id;", producing a single _struct_.
	* `SelectIds(struct, tableLocation)`: "SELECT ID FROM TABLE;", producing a slice of strings.
	* `SelectAll(struct, tableLocation)`: "SELECT * FROM TABLE;", producing a map[id]_struct_ (a map of structs, where the index is the table ID)
	* `SelectWhere(struct, function, tableLocation)`: "SELECT * FROM TABLE WHERE conditions;", producing a map[id]_struct_ (a map of structs, where the index is the table ID), according to a function, which can be a closure, a nested or a common function.

## Example usage

```
package main
import("fmt"; db "github.com/rodolfoap/jdocdb";)

// First, define the structure of your tables
type Person struct {
	Name string
	Age  int
	Sex  bool
}

type Animal struct {
	Name string
	Legs int
	Beak bool
}

// Now, you can start SQLing...
func main() {
	/* All functions have some PARAMETERS and then [ PREFIX [, SUFFIX] ],

	For example: db.SelectIds(Person{}, "prefix", "suffix")

	1. If NO PREFIX is specified, the path will be ./person/
	2. If PREFIX=data, the path will be ./data/person/
	3. If SUFFIX=people, the path will be ./data/people/

	*/

	/*: Usage: db.Insert(KEY, STRUCT, [ PREFIX [, SUFFIX] ]) */
	db.Insert("p0926", Person{"James", 33, false}, "prefix", "suffix")
	/* Will produce the following file: ./data/people/p0926.json:
	{
		"Id": "p0926",
		"Data": {
			"Name": "James",
			"Age": 33,
			"Sex": false
		}
	} */

	db.Insert("z0215", Person{"Jenna", 11, false})
	db.Insert("w1132", Person{"Joerg", 22, true}, "prefix")
	db.Insert("q9823", Person{"Jonas", 44, true}, "prefix", "suffix")
	db.Insert("r8791", Person{"Jonna", 55, false}, "prefix", "suffix")
	db.Insert("n9878", Person{"Junge", 55, true}, "prefix", "suffix")
	/* Now, we have:
	.
	├── person/
	│   └── z0215.json
	└── prefix/
	    ├── person/
	    │   └── w1132.json
	    └── suffix/
	        ├── n9878.json
	        ├── p0926.json
	        ├── q9823.json
	        └── r8791.json

	1. If NO PREFIX is specified, the path will be ./person/
	2. If PREFIX=data, the path will be ./data/person/
	3. If SUFFIX=people, the path will be ./data/people/
	*/

	/* Usage: db.Select(KEY, EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	jonas:=db.Select("q9823", Person{}, "prefix", "suffix")
	// {Jonas 44 true}, main.Person, 44
	fmt.Printf("%v, %T, %v\n", jonas, jonas, jonas.Age)

	/* Usage: db.SelectIds(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	listIds:=db.SelectIds(Person{}, "prefix", "suffix")
	// [n9878 p0926 q9823 r8791]
	fmt.Println(listIds)

	/* Usage: db.SelectAll(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	m:=db.SelectAll(Person{}, "prefix", "suffix")
	// map[n9878:{Junge 55 true} p0926:{James 33 false} q9823:{Jonas 44 true} r8791:{Jonna 55 false}]
	fmt.Println(m)

	/* Complex Queries over a Single Table **********************************************************/

	/* Do whatever query emulating a SELECT*FROM [TABLE] WHERE [CONDITIONS...] */
	/* Usage: db.SelectWhere(EMPTY_STRUCT, func(p Table) bool, [ PREFIX [, SUFFIX] ]) */

	filtered:=db.SelectWhere(Person{}, func(p Person) bool { return p.Age==55 }, "prefix", "suffix")
	// map[n9878:{Junge 55 true} r8791:{Jonna 55 false}]
	fmt.Println("Having 55:", filtered)

	filtered=db.SelectWhere(Person{}, func(p Person) bool { return !p.Sex }, "prefix", "suffix")
	// map[p0926:{James 33 false} r8791:{Jonna 55 false}]
	fmt.Println("Have not Sex:", filtered)

	filtered=db.SelectWhere(Person{}, func(p Person) bool { return p.Sex && p.Age==55 }, "prefix", "suffix")
	// map[n9878:{Junge 55 true}]
	fmt.Println("Have Sex and 55:", filtered)

	// Testing queries with a new table...
	db.Insert("dinosaur", Animal{"Barney", 2, false}, "prefix")
	db.Insert("chicken", Animal{"Clotilde", 2, true}, "prefix")
	db.Insert("dog", Animal{"Wallander, Mortimer", 4, false}, "prefix")
	db.Insert("cat", Animal{"Watson", 3, false}, "prefix")
	db.Insert("ant", Animal{"Woody", 5, true}, "prefix")

	/*
		Result:

		prefix/
		└── animal/
		    ├── ant.json
		    ├── cat.json
		    ├── chicken.json
		    ├── dinosaur.json
		    └── dog.json
	*/

	// A nested function, any kind of function will do.
	hasLongNameOrBeak:=func(a Animal) bool { return len(a.Name)>6 || a.Beak }

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	fmt.Println("Has Long Name Or Beak:", db.SelectWhere(Animal{}, hasLongNameOrBeak, "prefix"))
}
```
