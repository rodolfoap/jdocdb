# jdocdb: A JSON File-documents Database

This is a minimalist file-based JSON documents database. Tables are subdirectories, registries are files, filenames are registry IDs. That's it.

## TODO

* Needs to be thread-safe
* Needs error handling
* Needs some logging
* SELECT needs to be improved, the reflection loop is fragile and could fail under heavy conditions
* Needs better SELECT comparison operators, maybe passing functions

## Example usage

```
package main
import("fmt";   db "github.com/rodolfoap/jdocdb";)

type Person struct {
        Name string
        Age  int
        Sex  bool
}

func main() {

	/*: Usage: Insert(KEY, STRUCT) */
	db.Insert("p0926", Person{"James", 33, false})
	/* Will produce the following file: ./person/p0926.json:
	{
		"Id": "p0926",
		"Data": {
			"Name": "James",
			"Age": 33,
			"Sex": false
		}
	} */

	db.Insert("q9823", Person{"Jonas", 44, true})
	db.Insert("r8791", Person{"Jonna", 55, false})
	db.Insert("n9878", Person{"Junge", 55, true})
	/* Now, we have:
	.
	└── person
	    ├── n9878.json
	    ├── p0926.json
	    ├── q9823.json
	    └── r8791.json
	*/

	/* Usage: Select(KEY, EMPTY_STRUCT) */
	jonas:=db.Select("q9823", Person{})
	// {Jonas 44 true}, main.Person, 44
	fmt.Printf("%v, %T, %v\n", jonas, jonas, jonas.Age)

	/* Usage: SelectIds(EMPTY_STRUCT) */
	listIds:=db.SelectIds(Person{})
	// [n9878 p0926 q9823 r8791]
	fmt.Println(listIds)

	/* Usage: SelectAll(EMPTY_STRUCT) */
	m:=db.SelectAll(Person{})
	// map[n9878:{Junge 55 true} p0926:{James 33 false} q9823:{Jonas 44 true} r8791:{Jonna 55 false}]
	fmt.Println(m)

	/* Usage: SelectAll(EMPTY_STRUCT, strings conditions map: "Key": "Value") */
	/* Warning: Keys are case-sensitive. Only exact comparisons are available. */
	filtered:=db.SelectFilter(Person{}, map[string]string{"Age": "55"})
	// map[n9878:{Junge 55 true} r8791:{Jonna 55 false}]
	fmt.Println(filtered)

	filtered=db.SelectFilter(Person{}, map[string]string{"Sex": "false"})
	// map[p0926:{James 33 false} r8791:{Jonna 55 false}]
	fmt.Println(filtered)

	filtered=db.SelectFilter(Person{}, map[string]string{"Sex": "true", "Age": "55"})
	// map[n9878:{Junge 55 true}]
	fmt.Println(filtered)
}
```
