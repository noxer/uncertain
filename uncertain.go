package uncertain

import (
	"errors"
	"reflect"
	"strconv"
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
	if len(path) == 0 {
		return from, nil
	}

	for from.Kind() == reflect.Ptr || from.Kind() == reflect.Interface {
		if from.IsNil() {
			return reflect.Value{}, errors.New("unable to dereference nil pointer")
		}
		from = from.Elem()
	}

	switch from.Kind() {

	case reflect.Slice:
		if from.IsNil() {
			return reflect.Value{}, errors.New("unable to dereference nil slice")
		}
		fallthrough

	case reflect.Array, reflect.String:
		i, ok := anyToInt(path[0])
		if !ok || i < 0 || i >= from.Len() {
			return reflect.Value{}, errors.New("index is out of bounds")
		}
		return get(from.Index(i), path[1:])

	case reflect.Map:
		from = from.MapIndex(reflect.ValueOf(path[0]))
		if !from.IsValid() {
			return reflect.Value{}, errors.New("map key not found")
		}
		return get(from, path[1:])

	case reflect.Struct:
		field, ok := path[0].(string)
		if !ok {
			return reflect.Value{}, errors.New("invalid field name")
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
	}

	return 0, false
}
