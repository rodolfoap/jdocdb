package jdocdb
import("encoding/json"; "fmt"; "io/ioutil"; "os"; "path/filepath"; "reflect"; "strings";)

type Register struct {
	Id   string
	Data interface{}
}

func GetType(doc interface{}) string {
	return strings.ToLower(strings.SplitN(fmt.Sprintf("%T", doc), ".", 2)[1])
}

// Inserts one registry into a table using its ID
func Insert[T interface{}](id string, doc T) {
	reg:=Register{Id: id, Data: doc}
	table:=filepath.Clean(GetType(doc))
	os.MkdirAll(table, 0755)
	jsonPath:=filepath.Join(table, id+".json")
	jsonBytes, _:=json.MarshalIndent(reg, "", "\t")
	jsonBytes=append(jsonBytes, byte('\n'))
	ioutil.WriteFile(jsonPath, jsonBytes, 0644)
}

//Selects one registry from a table using its ID
func Select[T interface{}](id string, doc T) T {
	reg:=Register{Id: id, Data: &doc}
	table:=filepath.Clean(GetType(doc))
	jsonPath:=filepath.Join(table, id+".json")
	jsonBytes, _:=ioutil.ReadFile(jsonPath)
	json.Unmarshal(jsonBytes, &reg)
	return doc
}

// Select all IDs of a table
func SelectIds[T interface{}](doc T) []string {
	idList:=[]string{}
	table:=filepath.Clean(GetType(doc))
	fileList, _:=ioutil.ReadDir(table)
	for _, f:=range fileList {
		if strings.HasSuffix(f.Name(), ".json") {
			idList=append(idList, strings.TrimSuffix(f.Name(), ".json"))
		}
	}
	return idList
}

// Selects all rows from a table
func SelectAll[T interface{}](doc T) map[string]T {
	docs:=map[string]T{}
	for _, id:=range SelectIds(doc) {
		docs[id]=Select(id, doc)
	}
	return docs
}

// Cleans up strings for comparison
func neat(value interface{}) string {
	return strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", value)))
}

// Selects all rows that meet some conditions
func SelectFilter[T interface{}](doc T, cond map[string]string) map[string]T {
	docs:=map[string]T{}
	// Loop all documents of the table
	for id, one:=range SelectAll(doc) {
		accept:=true
		v:=reflect.ValueOf(one)
		// Now, loop the fields
		for i:=0; i<v.NumField(); i++ {
			fieldKey:=v.Type().Field(i).Name
			fieldValue:=v.Field(i).Interface()
			// Now, check if a condition of such fieldKey (e.g. cond["Age"]) exists
			// If so, exists=true and condValue=55
			if condValue, exists:=cond[fieldKey]; exists {
				if neat(fieldValue)!=neat(condValue) {
					accept=false
				}
			}
		}
		if accept { // If conditions have all passed (accept is still true)
			docs[id]=one
		}
	}
	return docs
}
