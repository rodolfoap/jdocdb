package jdocdb

import (
	"fmt"
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
	fmt.Printf("%v, %T, %v\n", jonas, jonas, jonas.Age)

	/* Usage: SelectIds(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	listIds := SelectIds(Person{}, "prefix", "suffix")
	// [n9878 p0926 q9823 r8791]
	fmt.Println(listIds)

	/* Usage: SelectAll(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	m := SelectAll(Person{}, "prefix", "suffix")
	// map[n9878:{Junge 55 true} p0926:{James 33 false} q9823:{Jonas 44 true} r8791:{Jonna 55 false}]
	fmt.Println(m)

	/* Complex Queries: do whatever query emulating a SELECT*FROM [TABLE] WHERE [CONDITIONS...] */
	/* Usage: SelectWhere(EMPTY_STRUCT, func(p Table) bool, [ PREFIX [, SUFFIX] ]) */
	/* Do not forget to declare the structure as a Table, see the top of this file */

	filtered := SelectWhere(Person{}, func(p Person) bool { return p.Age == 55 }, "prefix", "suffix")
	// map[n9878:{Junge 55 true} r8791:{Jonna 55 false}]
	fmt.Println("Having 55:", filtered)

	filtered = SelectWhere(Person{}, func(p Person) bool { return !p.Sex }, "prefix", "suffix")
	// map[p0926:{James 33 false} r8791:{Jonna 55 false}]
	fmt.Println("Have not Sex:", filtered)

	filtered = SelectWhere(Person{}, func(p Person) bool { return p.Sex && p.Age == 55 }, "prefix", "suffix")
	// map[n9878:{Junge 55 true}]
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

	// map[ant:{Woody 5 true} chicken:{Clotilde 2 true} dog:{Wallander, Mortimer 4 false}]
	fmt.Println("Has Long Name Or Beak:", SelectWhere(Animal{}, hasLongNameOrBeak, "prefix"))
}
