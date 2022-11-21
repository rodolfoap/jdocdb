package jdocdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
func Test_lib(t *testing.T) {
	/*

		All functions have some PARAMETERS and then [ PREFIX [, SUFFIX] ],

		For example: SelectIds(Person{}, "prefix", "suffix")

		1. If NO PREFIX is specified, the path will be ./person/
		2. If PREFIX=data, the path will be ./data/person/
		3. If SUFFIX=people, the path will be ./data/people/

		Example with Insert, usage:

		db.Insert(KEY, STRUCT, [ PREFIX [, SUFFIX] ])
	*/

	Insert("p0926", Person{"James", 33, false})
	/* Will create the file ./people/p0926.json with this content:
	{
		"Id": "p0926",
		"Data": {
			"Name": "James",
			"Age": 33,
			"Sex": false
		}
	} */

	Insert("dinosaur", Animal{"Barney", 2, false})
	// Will create ./animal/dinosaur.json

	/*

		Where is my table?

	*/

	// This creates /tmp/person/z0215.json (default suffix is the table name)
	// This is the recommended way to perform queries: just the database location (prefix).
	Insert("z0215", Person{"Junge", 19, true}, "/tmp")

	// This creates /tmp/z0215.json (no suffix)
	Insert("z0215", Person{"Junge", 19, true}, "/tmp", "")

	// When prefix and suffix are present: prefix/suffix/ID.json
	// Naturally, suffix can be absolute (/var/data/...) or relative (./table)
	// This will create ./prefix/suffix/z0215.json
	Insert("z0215", Person{"Junge", 11, true}, "prefix", "suffix")

	// When only prefix is present: prefix/TABLENAME/ID.json
	// Notice such is a different table with the same type
	// This will create ./prefix/person/z0215.json
	Insert("z0215", Person{"Junge", 19, true}, "prefix")

	// When none is present: ./TABLENAME/ID.json
	// This will create ./person/q9823.json
	Insert("q9823", Person{"Jonas", 44, true})

	// More data in the ./person/ table:
	Insert("n9878", Person{"Junge", 55, true})
	Insert("r8791", Person{"Jonna", 55, false})

	/* Excluding what was created in /tmp, here we have:
	.
	├── animal
	│   └── dinosaur.json
	├── prefix
	│   ├── person
	│   │   └── z0215.json
	│   └── suffix
	│       └── z0215.json
	└── person
	    ├── p0926.json
	    ├── q9823.json
	    ├── r8791.json
	    └── n9878.json

	1. If NO PREFIX is specified, the path will be ./person/
	2. If PREFIX=data, the path will be ./data/person/
	3. If SUFFIX=people, the path will be ./data/people/
	*/

	/* Usage: Select(KEY, EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	jonas := Select("q9823", Person{})

	// {Jonas 44 true}, jdocdb.Person, 44
	fmt.Printf("%v, %T, %v\n", jonas, jonas, jonas.Age)
	assert.Equal(t, jonas.Age, 44)
	assert.IsType(t, jonas, Person{})

	/* Usage: SelectIds(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	listIds := SelectIds(Person{})
	// [n9878 p0926 q9823 r8791]
	fmt.Println(listIds)
	assert.IsType(t, listIds, []string{})
	assert.Len(t, listIds, 4)
	assert.Contains(t, listIds, "n9878")

	/* Usage: SelectAll(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	m := SelectAll(Person{})
	// map[n9878:{Junge 55 true} p0926:{James 33 false} q9823:{Jonas 44 true} r8791:{Jonna 33 false}]
	fmt.Println(m)
	assert.IsType(t, m, map[string]Person{})
	assert.Len(t, m, 4)

	// Empty SELECT: file does not exist, and a log/debug message is produced
	jojo := Select("a7654", Person{})
	fmt.Printf("This is just empty: %v\n", jojo)

	/* COMPLEX QUERIES

	SELECT * FROM [TABLE] WHERE [CONDITIONS...]; are possible using Golang func types
	that are passed as arguments. See the following examples with SelectWhere:

	db.SelectWhere(EMPTY_STRUCT, func(p Table) bool, [ PREFIX [, SUFFIX] ])

	*/

	//
	// SELECT * FROM Person WHERE AGE == 55
	//
	filtered := SelectWhere(Person{}, func(p Person) bool { return p.Age == 55 })

	// map[n9878:{Junge 55 true} r8791:{Jonna 55 false}]
	assert.Len(t, filtered, 2)
	assert.Contains(t, filtered, "n9878", "r8791")
	fmt.Println("Having 55:", filtered)

	//
	// SELECT * FROM Person WHERE NOT Sex
	//
	filtered = SelectWhere(Person{}, func(p Person) bool { return !p.Sex })

	// map[p0926:{James 33 false} r8791:{Jonna 55 false}]
	fmt.Println("Have not Sex:", filtered)
	assert.Len(t, filtered, 2)
	assert.Contains(t, filtered, "p0926", "r8791")

	//
	// SELECT * FROM Person WHERE Sex AND AGE == 55
	//
	filtered = SelectWhere(Person{}, func(p Person) bool { return p.Sex && p.Age == 55 })

	// map[n9878:{Junge 55 true}]
	assert.Len(t, filtered, 1)
	assert.Contains(t, filtered, "n9878")
	fmt.Println("Have Sex and 55:", filtered)

	// Testing queries with a new table...
	Insert("chicken", Animal{"Clotilde", 2, true})
	Insert("dog", Animal{"Wallander, Mortimer", 4, false})
	Insert("cat", Animal{"Watson", 3, false})
	Insert("ant", Animal{"Woody", 5, true})

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
	animals := SelectWhere(Animal{}, hasLongNameOrBeak)

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	fmt.Println("Has Long Name Or Beak:", animals)
	assert.Len(t, animals, 3)
	assert.Contains(t, animals, "ant", "chicken", "dog")

	//
	// SELECT ID FROM Animal WHERE LEN(name)>6 OR Beak
	//
	animalIDs := SelectIdWhere(Animal{}, hasLongNameOrBeak)

	// [chicken dog ant]
	fmt.Println("IDs for Has Long Name Or Beak:", animalIDs)
	assert.Len(t, animals, 3)
	assert.Contains(t, animals, "ant", "chicken", "dog")

	/* QUERIES performing AGGREGATION

	SQL aggregation functions like SUM(), AVG() or COUNT() are possible in JDocDB in the
	same way that WHERE functions are possible: by passing anonymous or common functions as
	arguments.

	*/

	//
	// SELECT ... SUM(*) AS sum FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	sum := 0
	animals = SelectWhereAggreg(Animal{}, hasLongNameOrBeak, &sum, func(id string, a Animal) { sum += a.Legs })

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}],
	// sum == 11
	fmt.Printf("%v, have a total of %v Legs.\n", animals, sum)
	assert.Len(t, animals, 3)
	assert.Contains(t, animals, "ant", "chicken", "dog")
	assert.Equal(t, sum, 11)

	// Multiple aggregations
	//
	// SELECT ... COUNT(*) AS x0, SUM(Legs) AS x1 FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	x := []int{0, 0}
	animals = SelectWhereAggreg(Animal{}, hasLongNameOrBeak, &x, func(id string, a Animal) { x[0] += 1; x[1] += a.Legs })

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	// map[ant:{...} chicken:{...} dog:{...}], COUNT: 3; SUM(Legs): 11.
	fmt.Printf("%v, COUNT: %v; SUM(Legs): %v.\n", animals, x[0], x[1])
	assert.Len(t, animals, 3)
	assert.Contains(t, animals, "ant", "chicken", "dog")
	assert.Equal(t, x[0], 3)  // COUNT(*)
	assert.Equal(t, x[1], 11) // SUM(Legs)

	// Multiple aggregations without WHERE, use SelectAggreg()
	//
	// SELECT ... COUNT(*) AS x0, SUM(Legs) AS x1 FROM Animal;
	//
	x = []int{0, 0}
	animals = SelectAggreg(Animal{}, &x, func(id string, a Animal) { x[0] += 1; x[1] += a.Legs })

	// map[ant:{Woody 5 true} cat:{Watson 3 false} chicken:{Clotilde 2 true} dinosaur:{Barney 2 false} dog:{Wallander, Mortimer 4 false}]
	// map[ant:{...} chicken:{...} dog:{...}], COUNT: 3; SUM(Legs): 11.
	fmt.Printf("%v, COUNT: %v; SUM(Legs): %v.\n", animals, x[0], x[1])
	assert.Len(t, animals, 5)
	assert.Contains(t, animals, "ant", "cat", "chicken", "dinosaur", "dog")
	assert.Equal(t, x[0], 5)  // COUNT(*)
	assert.Equal(t, x[1], 16) // SUM(Legs)

	/* A simpler approach to aggregation: Count()

	If you need just aggregations, don't get the whole set, just get the count and
	get your aggregation result using Count().
	*/

	//
	// SELECT ... COUNT(*) AS legs FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	legs := 0
	quantity := CountWhereAggreg(Animal{}, hasLongNameOrBeak, &legs, func(id string, a Animal) { legs += a.Legs })

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	// quantity == 3 and legs == 11
	fmt.Printf("Simpler, COUNT: %v; SUM(Legs): %v.\n", quantity, legs)
	assert.Equal(t, quantity, 3)
	assert.Equal(t, legs, 11)

	// Simpler, without the WHERE clause:
	//
	// SELECT ... COUNT(*) AS legs FROM Animal;
	//
	legs = 0
	quantity = CountAggreg(Animal{}, &legs, func(id string, a Animal) { legs += a.Legs })
	// map[ant:{Woody 5 true} cat:{Watson 3 false} chicken:{Clotilde 2 true} dinosaur:{Barney 2 false} dog:{Wallander, Mortimer 4 false}]
	// quantity == 5 and legs == 16
	fmt.Printf("Even simpler, COUNT: %v; SUM(Legs): %v.\n", quantity, legs)
	assert.Equal(t, quantity, 5)
	assert.Equal(t, legs, 16)

	// Just Count() without aggregation, using a WHERE clause:
	//
	// SELECT ... COUNT(*) AS legs FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	quantity = CountWhere(Animal{}, hasLongNameOrBeak)

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	// quantity == 5 and legs == 16
	fmt.Printf("COUNT WHERE: %v.\n", quantity)
	assert.Equal(t, quantity, 3)

	// Simplest COUNT:
	//
	// SELECT ... COUNT(*) FROM Animal;
	//
	quantity = Count(Animal{})
	// map[ant:{Woody 5 true} cat:{Watson 3 false} chicken:{Clotilde 2 true} dinosaur:{Barney 2 false} dog:{Wallander, Mortimer 4 false}]
	// quantity == 5
	fmt.Printf("Bare COUNT: %v.\n", quantity)
	assert.Equal(t, quantity, 5)

	// Addition: SQL SUM is Sum()
	//
	// SELECT SUM(Legs) FROM Animal;
	//
	quantLegs := Sum(Animal{}, "Legs")
	assert.Equal(t, quantLegs, 16)
	fmt.Printf("SUM: %v.\n", quantLegs)

	// SUM WHERE
	//
	// SELECT SUM(Legs) FROM Animal WHERE LEN(name)>6 OR Beak;
	//
	quantLegs = SumWhere(Animal{}, "Legs", hasLongNameOrBeak)
	assert.Equal(t, quantLegs, 11)
	fmt.Printf("SUM WHERE: %v.\n", quantLegs)

	// Delete function
	//
	// DELETE FROM Person WHERE ID="p0926";
	//
	Delete("p0926", Person{})
	Delete("n9878", Person{})
	Delete("q9823", Person{})
	Delete("r8791", Person{})

	remaining := SelectAll(Person{})
	fmt.Println("Remaining after delete:", remaining)
	assert.Len(t, remaining, 0)
}
