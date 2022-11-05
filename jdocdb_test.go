package jdocdb
import("fmt"; "testing";)

type Person struct {
	Name string
	Age  int
	Sex  bool
}

func Test_lib(t *testing.T) {
	/*: Usage: Insert(KEY, STRUCT) */
	Insert("p0926", Person{"James", 33, false}, "data", "db")
	/* Will produce the following file: ./person/p0926.json:
	{
		"Id": "p0926",
		"Data": {
			"Name": "James",
			"Age": 33,
			"Sex": false
		}
	} */

	Insert("q9823", Person{"Jonas", 44, true}, "data", "db")
	Insert("r8791", Person{"Jonna", 55, false}, "data", "db")
	Insert("n9878", Person{"Junge", 55, true}, "data", "db")

	/* Usage: Select(KEY, EMPTY_STRUCT) */
	jonas:=Select("q9823", Person{}, "data", "db")
	// {Jonas 44 true}, main.Person, 44
	fmt.Printf("%v, %T, %v\n", jonas, jonas, jonas.Age)

	/* Usage: SelectIds(EMPTY_STRUCT) */
	listIds:=SelectIds(Person{}, "data", "db")
	// [n9878 p0926 q9823 r8791]
	fmt.Println(listIds)

	/* Usage: SelectAll(EMPTY_STRUCT) */
	m:=SelectAll(Person{}, "data", "db")
	// map[n9878:{Junge 55 true} p0926:{James 33 false} q9823:{Jonas 44 true} r8791:{Jonna 55 false}]
	fmt.Println(m)

	/* Usage: SelectAll(EMPTY_STRUCT, strings conditions map: "Key": "Value") */
	/* Warning: Keys are case-sensitive. Only exact comparisons are available. */
	filtered:=SelectFilter(Person{}, map[string]string{"Age": "55"}, "data", "db")
	// map[n9878:{Junge 55 true} r8791:{Jonna 55 false}]
	fmt.Println(filtered)

	filtered=SelectFilter(Person{}, map[string]string{"Sex": "false"}, "data", "db")
	// map[p0926:{James 33 false} r8791:{Jonna 55 false}]
	fmt.Println(filtered)

	filtered=SelectFilter(Person{}, map[string]string{"Sex": "true", "Age": "55"}, "data", "db")
	// map[n9878:{Junge 55 true}]
	fmt.Println(filtered)

}
