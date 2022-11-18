package jdocdb
import("testing"; b "github.com/rodolfoap/gx";)

func SelectWhere[T Table](doc T, cond func(T) bool, prefix ...string) bool {
	return cond(doc)
}

type Person struct {
	Name string
	Age  int
	Sex  bool
}

type Animal struct {
	Name  string
	Legs  int
	Wings bool
}

type Table interface {
	Person | Animal
}

func Test_select(t *testing.T) {
	b.Trace("HASSEX:", SelectWhere(Person{Sex: true}, func(p Person) bool { return p.Sex }, "prefix", "suffix")) //TRUE
}
