# jdocdb: A JSON File-documents Database
[![Go Reference](https://pkg.go.dev/badge/github.com/rodolfoap/jdocdb.svg)](https://pkg.go.dev/github.com/rodolfoap/jdocdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/rodolfoap/jdocdb)](https://goreportcard.com/report/github.com/rodolfoap/jdocdb)
[![Coverage Status](https://coveralls.io/repos/github/rodolfoap/jdocdb/badge.svg?branch=main)](https://coveralls.io/github/rodolfoap/jdocdb?branch=main)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
[![GPLv3 license](https://img.shields.io/badge/License-GPLv3-blue.svg)](http://perso.crans.org/besson/LICENSE.html)

A minimalist file-based JSON documents database with the capability of complex SELECT WHERE operations and multiple aggregations over a single table.

* Tables are subdirectories, e.g. `/tmp/client/`;
* Registries are files, e.g. `/tmp/client/a929782.json` that map to golang _structs_.
* Filenames are registry IDs, e.g. `/tmp/client/a929782.json` has `ID==a929782`;
* SQL queries can include WHERE clauses (e.g. `SelectWhere()`), aggregations (`SelectAggreg()`, to perform SUM(), COUNT(), AVG(), advanced filters or far any complex result) or both (`SelectWhereAggreg()`).

## Install

```
go get github.com/rodolfoap/jdocdb
```

## JDocDB HelloWorld

```
package main
import("fmt"; db "github.com/rodolfoap/jdocdb";)

type Client struct {
        Name   string
        Age    int
        Active bool
}

func main() {
        // Insert
        c:=Client{"Hello, World!", 44, true}
        db.Insert("a929782", c, "/tmp")
        /*
        cat /tmp/client/a929782.json
        {
                "Id": "a929782",
                "Data": {
                        "Name": "Hello, World!",
                        "Age": 44,
                        "Active": true
                }
        }
        */

        // Select
        result1:=db.Select("a929782", Client{}, "/tmp")
        fmt.Printf("%#v\n", result1) // main.Client{Name:"Hello, World!", Age:44, Active:true}
}
```
## Equivalencies with SQL

JDocDB has far stronger query potentials in addition to SELECT by ID. `SQL SELECT` equivalents are:

* `SELECT * FROM mytype;`: SelectAll(MyType{}, tableLocation), producing a **map[id]_struct_** (a map of _MyType_ structs, where the index is the register ID)
* `SELECT ID FROM mytype;`: `SelectIds(MyType{}, tableLocation)`, producing a slice of string IDs.
* `SELECT * FROM mytype WHERE ID=id;`: `Select(id, MyType{}, tableLocation)`, producing a single _struct_, the ID is necessarily unique.
* `SELECT * FROM mytype WHERE conditions...;`: `SelectWhere(MyType{}, conditionsFunction, tableLocation)`, producing a map[id]_struct_ (a map of _MyType_ structs, where the index is the table ID), according to a function, which can be a closure, a nested or a common function.
* `SELECT some_aggregate_function(...) FROM mytype;"`: `SelectAggreg(MyType, &aggregationVariable, aggregationFunction, tableLocation)`, producing a map[id]_struct_ (a map of _MyType_ structs, where the index is the table ID), according to a function, which can be a closure, a nested or a common function.
* `SELECT some_aggregate_function(...) FROM mytype WHERE conditions ...;"`: `SelectWhereAggreg(MyType, function, &aggregationVariable, aggregationFunction, tableLocation)`, producing a map[id]_struct_ (a map of _MyType_ structs, where the index is the table ID), according to a function, which can be a closure, a nested or a common function.
* `SELECT ID FROM mytype WHERE conditions...;`: `SelectIdWhere(MyType{}, conditionsFunction, tableLocation)`, producing a slice of string IDs, according to a function, which can be a closure, a nested or a common function.
* `SELECT SUM(myfield) FROM mytype;`: `Sum(MyType{}, "myfield", tableLocation)`, producing a slice of string IDs, using a string to locate the field to sum (evidently, a single sum). This function can be easily implemented with `SelectAggreg()` or `SelectWhereAggreg()` and it is provided just for simplicity.
* `SELECT SUM(myfield) FROM mytype WHERE conditions...;`: `SumWhere(MyType{}, conditionsFunction, tableLocation)`, producing a slice of string IDs, according to a function, which can be a closure, a nested or a common function. Notice there's no `SumWhereAggreg()` function, since SUM is already an aggregation.
* `SELECT COUNT(*) FROM mytype;`: `db.Count(MyType{}, tableLocation)`, yielding an **int**. Provided for simplicity.
* `SELECT COUNT(*) FROM mytype WHERE conditions...;`: `db.CountWhere(MyType{}, conditionsFunction, tableLocation)`, yielding an **int**, see `SelectWhere()`.
* `SELECT COUNT(*), some_aggregate_function(...) FROM mytype;`: `db.CountAggreg(MyType{}, &aggregationVariable, aggregationFunction, tableLocation)`, yielding an **int**, see `SelectAggreg()`. If only aggregation results are required, use `db.CountAggreg()` instead of `SelectAggreg()`.
* `SELECT COUNT(*), some_aggregate_function(...) FROM mytype WHERE conditions...;`: `db.CountWhereAggreg(MyType{}, conditionsFunction, tableLocation)`, yielding an **int**, see `SelectWhereAggreg()`.
* `DELETE FROM mytype WHERE ID=id;`: `db.Delete(MyType{}, tableLocation)`, simple deletion. There is no table delete, just remove the directory where the registers are.
* `INSERT INTO mytype VALUES ...;`: `db.Insert(MyType{}, tableLocation)`, simple _upsert_ function.

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
	/*

		All functions have some PARAMETERS and then [ PREFIX [, SUFFIX] ],

		For example: db.SelectIds(Person{}, "prefix", "suffix")

		1. If NO PREFIX is specified, the path will be ./person/
		2. If PREFIX=data, the path will be ./data/person/
		3. If SUFFIX=people, the path will be ./data/people/

		Example with db.Insert, usage:

		db.Insert(KEY, STRUCT, [ PREFIX [, SUFFIX] ])
	*/

	db.Insert("p0926", Person{"James", 33, false})
	/* Will create the file ./people/p0926.json with this content:
	{
		"Id": "p0926",
		"Data": {
			"Name": "James",
			"Age": 33,
			"Sex": false
		}
	} */

	db.Insert("dinosaur", Animal{"Barney", 2, false})
	// Will create ./animal/dinosaur.json

	/*

		Where is my table?

	*/

	// This creates /tmp/person/z0215.json (default suffix is the table name)
	// This is the recommended way to perform queries: just the database location (prefix).
	db.Insert("z0215", Person{"Junge", 19, true}, "/tmp")

	// This creates /tmp/z0215.json (no suffix)
	db.Insert("z0215", Person{"Junge", 19, true}, "/tmp", "")

	// When prefix and suffix are present: prefix/suffix/ID.json
	// Naturally, suffix can be absolute (/var/data/...) or relative (./table)
	// This will create ./prefix/suffix/z0215.json
	db.Insert("z0215", Person{"Junge", 11, true}, "prefix", "suffix")

	// When only prefix is present: prefix/TABLENAME/ID.json
	// Notice such is a different table with the same type
	// This will create ./prefix/person/z0215.json
	db.Insert("z0215", Person{"Junge", 19, true}, "prefix")

	// When none is present: ./TABLENAME/ID.json
	// This will create ./person/q9823.json
	db.Insert("q9823", Person{"Jonas", 44, true})

	// More data in the ./person/ table:
	db.Insert("n9878", Person{"Junge", 55, true})
	db.Insert("r8791", Person{"Jonna", 55, false})

	/* Excluding what was created in /tmp, here we have:
	.
	├── animal
	│   └── dinosaur.json
	├── prefix
	│   ├── person
	│   │   └── z0215.json
	│   └── suffix
	│       └── z0215.json
	└── person
	    ├── p0926.json
	    ├── q9823.json
	    ├── r8791.json
	    └── n9878.json

	1. If NO PREFIX is specified, the path will be ./person/
	2. If PREFIX=data, the path will be ./data/person/
	3. If SUFFIX=people, the path will be ./data/people/
	*/

	/* Usage: db.Select(KEY, EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	jonas := db.Select("q9823", Person{})

	// {Jonas 44 true}, jdocdb.Person, 44
	fmt.Printf("%v, %T, %v\n", jonas, jonas, jonas.Age)

	/* Usage: db.SelectIds(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	listIds := db.SelectIds(Person{})
	// [n9878 p0926 q9823 r8791]
	fmt.Println(listIds)

	/* Usage: db.SelectAll(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	m := db.SelectAll(Person{})
	// map[n9878:{Junge 55 true} p0926:{James 33 false} q9823:{Jonas 44 true} r8791:{Jonna 33 false}]
	fmt.Println(m)

	// Empty SELECT: file does not exist, and a log/debug message is produced
	jojo := db.Select("a7654", Person{})
	fmt.Printf("This is just empty: %v\n", jojo)

	/* COMPLEX QUERIES

	SELECT * FROM [TABLE] WHERE [CONDITIONS...]; are possible using Golang func types
	that are passed as arguments. See the following examples with db.SelectWhere:

	db.SelectWhere(EMPTY_STRUCT, func(p Table) bool, [ PREFIX [, SUFFIX] ])

	*/

	//
	// SELECT * FROM Person WHERE AGE == 55
	//
	filtered := db.SelectWhere(Person{}, func(p Person) bool { return p.Age == 55 })

	// map[n9878:{Junge 55 true} r8791:{Jonna 55 false}]
	fmt.Println("Having 55:", filtered)

	//
	// SELECT * FROM Person WHERE NOT Sex
	//
	filtered = db.SelectWhere(Person{}, func(p Person) bool { return !p.Sex })

	// map[p0926:{James 33 false} r8791:{Jonna 55 false}]
	fmt.Println("Have not Sex:", filtered)

	//
	// SELECT * FROM Person WHERE Sex AND AGE == 55
	//
	filtered = db.SelectWhere(Person{}, func(p Person) bool { return p.Sex && p.Age == 55 })

	// map[n9878:{Junge 55 true}]
	fmt.Println("Have Sex and 55:", filtered)

	// Testing queries with a new table...
	db.Insert("chicken", Animal{"Clotilde", 2, true})
	db.Insert("dog", Animal{"Wallander, Mortimer", 4, false})
	db.Insert("cat", Animal{"Watson", 3, false})
	db.Insert("ant", Animal{"Woody", 5, true})

	/* Result:
	.
	└── animal
	    ├── ant.json
	    ├── cat.json
	    ├── chicken.json
	    ├── dinosaur.json
	    └── dog.json
	*/

	// An alternative way of Golang for defining functions:
	hasLongNameOrBeak := func(a Animal) bool { return len(a.Name) > 6 || a.Beak }

	//
	// SELECT * FROM Animal WHERE LEN(name)>6 OR Beak
	//
	animals := db.SelectWhere(Animal{}, hasLongNameOrBeak)

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	fmt.Println("Has Long Name Or Beak:", animals)

	//
	// SELECT ID FROM Animal WHERE LEN(name)>6 OR Beak
	//
	animalIDs := db.SelectIdWhere(Animal{}, hasLongNameOrBeak)

	// [chicken dog ant]
	fmt.Println("IDs for Has Long Name Or Beak:", animalIDs)

	/* QUERIES performing AGGREGATION

	SQL aggregation functions like SUM(), AVG() or COUNT() are possible in JDocDB in the
	same way that WHERE functions are possible: by passing anonymous or common functions as
	arguments.

	*/

	//
	// SELECT ... SUM(*) AS sum FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	sum := 0
	animals = db.SelectWhereAggreg(Animal{}, hasLongNameOrBeak, &sum, func(id string, a Animal) { sum += a.Legs })

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}],
	// sum == 11
	fmt.Printf("%v, have a total of %v Legs.\n", animals, sum)

	// Multiple aggregations
	//
	// SELECT ... COUNT(*) AS x0, SUM(Legs) AS x1 FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	x := []int{0, 0}
	animals = db.SelectWhereAggreg(Animal{}, hasLongNameOrBeak, &x, func(id string, a Animal) { x[0] += 1; x[1] += a.Legs })

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	// map[ant:{...} chicken:{...} dog:{...}], COUNT: 3; SUM(Legs): 11.
	fmt.Printf("%v, COUNT: %v; SUM(Legs): %v.\n", animals, x[0], x[1])

	// Multiple aggregations without WHERE, use db.SelectAggreg()
	//
	// SELECT ... COUNT(*) AS x0, SUM(Legs) AS x1 FROM Animal;
	//
	x = []int{0, 0}
	animals = db.SelectAggreg(Animal{}, &x, func(id string, a Animal) { x[0] += 1; x[1] += a.Legs })

	// map[ant:{Woody 5 true} cat:{Watson 3 false} chicken:{Clotilde 2 true} dinosaur:{Barney 2 false} dog:{Wallander, Mortimer 4 false}]
	// map[ant:{...} chicken:{...} dog:{...}], COUNT: 3; SUM(Legs): 11.
	fmt.Printf("%v, COUNT: %v; SUM(Legs): %v.\n", animals, x[0], x[1])

	/* A simpler approach to aggregation: db.Count()

	If you need just aggregations, don't get the whole set, just get the count and
	get your aggregation result using db.Count().
	*/

	//
	// SELECT ... COUNT(*) AS legs FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	legs := 0
	quantity := db.CountWhereAggreg(Animal{}, hasLongNameOrBeak, &legs, func(id string, a Animal) { legs += a.Legs })

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	// quantity == 3 and legs == 11
	fmt.Printf("Simpler, COUNT: %v; SUM(Legs): %v.\n", quantity, legs)

	// Simpler, without the WHERE clause:
	//
	// SELECT ... COUNT(*) AS legs FROM Animal;
	//
	legs = 0
	quantity = db.CountAggreg(Animal{}, &legs, func(id string, a Animal) { legs += a.Legs })
	// map[ant:{Woody 5 true} cat:{Watson 3 false} chicken:{Clotilde 2 true} dinosaur:{Barney 2 false} dog:{Wallander, Mortimer 4 false}]
	// quantity == 5 and legs == 16
	fmt.Printf("Even simpler, COUNT: %v; SUM(Legs): %v.\n", quantity, legs)

	// Just db.Count() without aggregation, using a WHERE clause:
	//
	// SELECT ... COUNT(*) AS legs FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	quantity = db.CountWhere(Animal{}, hasLongNameOrBeak)

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	// quantity == 5 and legs == 16
	fmt.Printf("COUNT WHERE: %v.\n", quantity)

	// Simplest COUNT:
	//
	// SELECT ... COUNT(*) FROM Animal;
	//
	quantity = db.Count(Animal{})
	// map[ant:{Woody 5 true} cat:{Watson 3 false} chicken:{Clotilde 2 true} dinosaur:{Barney 2 false} dog:{Wallander, Mortimer 4 false}]
	// quantity == 5
	fmt.Printf("Bare COUNT: %v.\n", quantity)

	// Addition: SQL SUM is db.Sum()
	//
	// SELECT SUM(Legs) FROM Animal;
	//
	quantLegs := db.Sum(Animal{}, "Legs")
	fmt.Printf("SUM: %v.\n", quantLegs)

	// SUM WHERE
	//
	// SELECT SUM(Legs) FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	quantLegs = db.SumWhere(Animal{}, "Legs", hasLongNameOrBeak)
	fmt.Printf("SUM WHERE: %v.\n", quantLegs)

	// db.Delete function
	//
	// DELETE FROM Person WHERE ID="p0926";
	//
	db.Delete("p0926", Person{})
	db.Delete("n9878", Person{})
	db.Delete("q9823", Person{})
	db.Delete("r8791", Person{})

	remaining := db.SelectAll(Person{})
	fmt.Println("Remaining after delete:", remaining)
}
```
