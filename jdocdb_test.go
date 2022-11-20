package jdocdb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
	/* All functions have some PARAMETERS and then [ PREFIX [, SUFFIX] ],

	For example: SelectIds(Person{}, "prefix", "suffix")

	1. If NO PREFIX is specified, the path will be ./person/
	2. If PREFIX=data, the path will be ./data/person/
	3. If SUFFIX=people, the path will be ./data/people/

	*/

	/*: Usage: Insert(KEY, STRUCT, [ PREFIX [, SUFFIX] ]) */
	Insert("p0926", Person{"James", 33, false}, "prefix", "suffix")
	/* Will produce the following file: ./data/people/p0926.json:
	{
		"Id": "p0926",
		"Data": {
			"Name": "James",
			"Age": 33,
			"Sex": false
		}
	} */

	Insert("z0215", Person{"Jenna", 11, false})
	Insert("w1132", Person{"Joerg", 22, true}, "prefix")
	Insert("q9823", Person{"Jonas", 44, true}, "prefix", "suffix")
	Insert("r8791", Person{"Jonna", 55, false}, "prefix", "suffix")
	Insert("n9878", Person{"Junge", 55, true}, "prefix", "suffix")
	/* Now, we have:
	.
	├── person
	│   └── z0215.json
	└── prefix
	    ├── person
	    │   └── w1132.json
	    └── suffix
	        ├── n9878.json
	        ├── p0926.json
	        ├── q9823.json
	        └── r8791.json

	1. If NO PREFIX is specified, the path will be ./person/
	2. If PREFIX=data, the path will be ./data/person/
	3. If SUFFIX=people, the path will be ./data/people/
	*/

	/* Usage: Select(KEY, EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	jonas := Select("q9823", Person{}, "prefix", "suffix")
	// {Jonas 44 true}, main.Person, 44
	assert.Equal(t, jonas.Age, 44)
	assert.IsType(t, jonas, Person{})
	fmt.Printf("%v, %T, %v\n", jonas, jonas, jonas.Age)

	/* Usage: SelectIds(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	listIds := SelectIds(Person{}, "prefix", "suffix")
	// [n9878 p0926 q9823 r8791]
	assert.IsType(t, listIds, []string{})
	assert.Len(t, listIds, 4)
	assert.Contains(t, listIds, "n9878")
	fmt.Println(listIds)

	/* Usage: SelectAll(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	m := SelectAll(Person{}, "prefix", "suffix")
	// map[n9878:{Junge 55 true} p0926:{James 33 false} q9823:{Jonas 44 true} r8791:{Jonna 55 false}]
	assert.IsType(t, m, map[string]Person{})
	assert.Len(t, m, 4)
	fmt.Println(m)

	// A bad SELECT: file does not exist
	jojo := Select("a7654", Person{}, "prefix", "suffix")
	fmt.Printf("This is just empty: %v\n", jojo)

	/* Complex Queries: do whatever query emulating a SELECT*FROM [TABLE] WHERE [CONDITIONS...] */
	/* Usage: SelectWhere(EMPTY_STRUCT, func(p Table) bool, [ PREFIX [, SUFFIX] ]) */
	/* Do not forget to declare the structure as a Table, see the top of this file */

	/*
		SELECT * FROM Person WHERE AGE == 55
	*/
	filtered := SelectWhere(Person{}, func(p Person) bool { return p.Age == 55 }, "prefix", "suffix")
	// map[n9878:{Junge 55 true} r8791:{Jonna 55 false}]
	assert.Len(t, filtered, 2)
	assert.Contains(t, filtered, "n9878", "r8791")
	fmt.Println("Having 55:", filtered)

	/*
		SELECT * FROM Person WHERE NOT Sex
	*/
	filtered = SelectWhere(Person{}, func(p Person) bool { return !p.Sex }, "prefix", "suffix")
	// map[p0926:{James 33 false} r8791:{Jonna 55 false}]
	assert.Len(t, filtered, 2)
	assert.Contains(t, filtered, "p0926", "r8791")
	fmt.Println("Have not Sex:", filtered)

	/*
		SELECT * FROM Person WHERE Sex AND AGE == 55
	*/
	filtered = SelectWhere(Person{}, func(p Person) bool { return p.Sex && p.Age == 55 }, "prefix", "suffix")
	// map[n9878:{Junge 55 true}]
	assert.Len(t, filtered, 1)
	assert.Contains(t, filtered, "n9878")
	fmt.Println("Have Sex and 55:", filtered)

	// Testing queries with a new table...
	Insert("dinosaur", Animal{"Barney", 2, false}, "prefix")
	Insert("chicken", Animal{"Clotilde", 2, true}, "prefix")
	Insert("dog", Animal{"Wallander, Mortimer", 4, false}, "prefix")
	Insert("cat", Animal{"Watson", 3, false}, "prefix")
	Insert("ant", Animal{"Woody", 5, true}, "prefix")

	/*
		Result:

		prefix/
		└── animal
		    ├── ant.json
		    ├── cat.json
		    ├── chicken.json
		    ├── dinosaur.json
		    └── dog.json
	*/

	// A nested function, any kind of function will do.
	hasLongNameOrBeak := func(a Animal) bool { return len(a.Name) > 6 || a.Beak }

	/*
		Example SELECT * WHERE LEN(name)>6 OR Beak
	*/
	animals := SelectWhere(Animal{}, hasLongNameOrBeak, "prefix")
	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	assert.Len(t, animals, 3)
	assert.Contains(t, animals, "ant", "chicken", "dog")
	fmt.Println("Has Long Name Or Beak:", animals)

	/*
		Example SELECT ID WHERE LEN(name)>6 OR Beak
	*/
	animalIDs := SelectIdWhere(Animal{}, hasLongNameOrBeak, "prefix")
	// [chicken dog ant]
	assert.Len(t, animals, 3)
	assert.Contains(t, animals, "ant", "chicken", "dog")
	fmt.Println("IDs for Has Long Name Or Beak:", animalIDs)

	/*
		Making a single aggregation, example: SELECT ... COUNT(*) AS sum
	*/
	sum := 0
	animals = SelectWhereGroup(Animal{}, hasLongNameOrBeak, &sum, func(a Animal) { sum += a.Legs }, "prefix")
	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	// sum == 11
	assert.Len(t, animals, 3)
	assert.Contains(t, animals, "ant", "chicken", "dog")
	assert.Equal(t, sum, 11)
	fmt.Printf("%v, have a total of %v Legs.\n", animals, sum)

	/*
		Making multiple aggregations, example: SELECT ... COUNT(*) AS x0, SUM(Legs) AS x1
	*/
	x := []int{0, 0}
	animals = SelectWhereGroup(Animal{}, hasLongNameOrBeak, &x, func(a Animal) { x[0] += 1; x[1] += a.Legs }, "prefix")
	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	// sum == 11
	assert.Len(t, animals, 3)
	assert.Contains(t, animals, "ant", "chicken", "dog")
	assert.Equal(t, x[0], 3)  // COUNT(*)
	assert.Equal(t, x[1], 11) // SUM(Legs)
	fmt.Printf("%v, COUNT: %v; SUM(Legs): %v.\n", animals, x[0], x[1])
	// map[ant:{...} chicken:{...} dog:{...}], COUNT: 3; SUM(Legs): 11.
}
