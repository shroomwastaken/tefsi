package domain

import (
	"fmt"
	"reflect"
	"strings"
)

// TODO: filter by price
type Filter struct {
	CategoryID   int
	SearchString string
}

// returns a slice of strings corresponding to each nonempty filter parameter,
// replacing sql parameters with $$$
//
// for example if filter.CategoryID == 2 and everything else is the zero value,
// this will return []string{ "items.category = $$$" }
// TODO: not hardcode table and column names
func (f *Filter) getNonEmptyParams() []string {
	result := []string{}

	if f.CategoryID != 0 {
		result = append(result, "items.category = $$$")
	}

	if f.SearchString != "" {
		// search matches searchString with title or description of the item
		result = append(result, "items.title LIKE CONCAT('%', CAST($$$ AS text), '%')")
	}

	return result
}

// generates string to append to sql query
func (f *Filter) GenerateString() string {
	filterParams := f.getNonEmptyParams()

	if len(filterParams) == 0 {
		return ""
	}

	str := "\nWHERE\n"

	for i, param := range filterParams {
		str += param
		if i < len(filterParams)-1 {
			str += "AND\n"
		}
	}

	cur := 1
	count := strings.Count(str, "$$$")
	for i := 0; i < count; i += 1 {
		str = strings.Replace(str, "$$$", fmt.Sprintf("$%v", cur), 1)
		cur += 1
	}

	return str
}

// returns parameters to pass to a db.Query() function
// should be used as db.Query(context, sqlString, filter.GetDBParams()...)
//
// only needs a list in the same order as the fields so uses silly reflect stuff
// to iterate over filters fields in a loop and compare to their zero value
func (f *Filter) GetDBParams() []any {
	result := []any{}

	filterValue := reflect.ValueOf(*f)

	for i := 0; i < filterValue.NumField(); i += 1 {
		curField := filterValue.Field(i).Interface()

		// skip if equal to zero value
		if curField != reflect.Zero(filterValue.Field(i).Type()).Interface() {
			result = append(result, curField)
		}
	}

	return result
}
