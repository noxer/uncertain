package uncertain

import (
	"errors"
	"reflect"
	"strconv"
	"unicode"
)

// Get the data nested in from when following path. You can access maps by providing the key,
// arrays, slices and strings by the index and structs by the field name. Pointers and interfaces
// are dereferenced except when they are the last item in the path.
func Get(from interface{}, path ...interface{}) (interface{}, error) {
	v, err := get(reflect.ValueOf(from), path)
	if err != nil {
		return nil, err
	}
	if !v.IsValid() {
		return nil, nil
	}
	return v.Interface(), nil
}

func get(from reflect.Value, path []interface{}) (reflect.Value, error) {
	// no more path segments, this seems to be the value the user wants
	if len(path) == 0 {
		return from, nil
	}

	// dereference all (nested) pointers and interfaces
	for from.Kind() == reflect.Ptr || from.Kind() == reflect.Interface {
		if from.IsNil() {
			return reflect.Value{}, errors.New("unable to dereference nil pointer")
		}
		from = from.Elem()
	}

	// now we try all the different (supported) kinds of values
	switch from.Kind() {

	case reflect.Slice:
		// this is a slice, check if it's nil before proceeding
		if from.IsNil() {
			return reflect.Value{}, errors.New("unable to dereference nil slice")
		}
		fallthrough

	case reflect.Array, reflect.String:
		// this is a slice, array, or string; use the path segment as an index
		i, ok := anyToInt(path[0])
		if !ok || i < 0 || i >= from.Len() {
			return reflect.Value{}, errors.New("index is out of bounds")
		}
		return get(from.Index(i), path[1:])

	case reflect.Map:
		// this is a map; use the path segment as the key
		from = from.MapIndex(reflect.ValueOf(path[0]))
		if !from.IsValid() {
			return reflect.Value{}, errors.New("map key not found")
		}
		return get(from, path[1:])

	case reflect.Struct:
		// this is a struct; use the path segment as a field name
		field, ok := path[0].(string)
		if !ok {
			return reflect.Value{}, errors.New("invalid field name")
		}
		if len(field) > 0 && !unicode.IsUpper([]rune(field)[0]) {
			return reflect.Value{}, errors.New("can't access private field")
		}
		from = from.FieldByName(field)
		if !from.IsValid() {
			return reflect.Value{}, errors.New("field not found in struct")
		}
		return get(from, path[1:])
	}

	return reflect.Value{}, errors.New("can't walk the rest of the path for" + from.String())
}

func anyToInt(i interface{}) (int, bool) {
	switch t := i.(type) {
	case string:
		n, err := strconv.Atoi(t)
		return n, err == nil
	case int:
		return t, true
	case uint:
		return int(t), true
	case int8:
		return int(t), true
	case uint8:
		return int(t), true
	case int16:
		return int(t), true
	case uint16:
		return int(t), true
	case int32:
		return int(t), true
	case uint32:
		return int(t), true
	case int64:
		return int(t), true
	case uint64:
		return int(t), true
	case uintptr:
		return int(t), true
	}

	return 0, false
}
