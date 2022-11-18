package jdocdb
import("fmt"; "testing"; b "github.com/rodolfoap/gx";)

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
	jonas:=Select("q9823", Person{}, "prefix", "suffix")
	// {Jonas 44 true}, main.Person, 44
	fmt.Printf("%v, %T, %v\n", jonas, jonas, jonas.Age)

	/* Usage: SelectIds(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	listIds:=SelectIds(Person{}, "prefix", "suffix")
	// [n9878 p0926 q9823 r8791]
	fmt.Println(listIds)

	/* Usage: SelectAll(EMPTY_STRUCT, [ PREFIX [, SUFFIX] ]) */
	m:=SelectAll(Person{}, "prefix", "suffix")
	// map[n9878:{Junge 55 true} p0926:{James 33 false} q9823:{Jonas 44 true} r8791:{Jonna 55 false}]
	fmt.Println(m)

	/* Usage: SelectAll(EMPTY_STRUCT, strings conditions map: "Key": "Value", [ PREFIX [, SUFFIX] ]) */
	/* Warning: Keys are case-sensitive. Only exact comparisons are available. */
	filtered:=SelectFilter(Person{}, map[string]string{"Age": "55"}, "prefix", "suffix")
	// map[n9878:{Junge 55 true} r8791:{Jonna 55 false}]
	fmt.Println(filtered)

	filtered=SelectFilter(Person{}, map[string]string{"Sex": "false"}, "prefix", "suffix")
	// map[p0926:{James 33 false} r8791:{Jonna 55 false}]
	fmt.Println(filtered)

	filtered=SelectFilter(Person{}, map[string]string{"Sex": "true", "Age": "55"}, "prefix", "suffix")
	// map[n9878:{Junge 55 true}]
	fmt.Println(filtered)

	check:=SelectFilter2(Person{}, hasSex, "prefix", "suffix")
	fmt.Println("HASSEX:", check)

	fmt.Println("CHECK:", Person{Sex: true}.CheckSex())
	fmt.Println("CHECK:", SelectFilter3(Person{Sex: true}.CheckSex()))
}

type fPerson func(person Person) bool

func hasSex(person Person) bool {
	return person.Sex
}
func SelectFilter2(doc Person, cond fPerson, prefix ...string) bool {
	if cond(doc) {
		b.Trace("TRUE")
		return true
	} else {
		b.Trace("FALSE")
		return false
	}
}
func SelectFilter3(function fPerson, prefix ...string) bool {
	if doc.function() {
		b.Trace("TRUE")
		return true
	} else {
		b.Trace("FALSE")
		return false
	}
}

/*===================================================================================*/
type Person struct {
	Name string
	Age  int
	Sex  bool
}
type Checker interface {
	Check() bool
}

func (p Person) CheckSex() bool {
	return p.Sex
}
