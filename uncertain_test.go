package uncertain_test

import (
	"fmt"
	"testing"

	"github.com/carlmjohnson/be"
	"github.com/noxer/uncertain"
)

type testStruct struct {
	Attr1 string
	Attr2 int
	Attr3 map[string]string
	Attr4 interface{}
}

func TestInvalid(t *testing.T) {
	_, err := uncertain.Get(func() {}, "segment")
	be.Nonzero(t, err)

	_, err = uncertain.Get(make(chan string), "segment")
	be.Nonzero(t, err)
}

func TestNil(t *testing.T) {
	res, err := uncertain.Get(nil)
	be.NilErr(t, err)
	be.Zero(t, res)

	_, err = uncertain.Get(nil, "segment")
	be.Nonzero(t, err)
	_, err = uncertain.Get(nil, 1)
	be.Nonzero(t, err)
}

func TestSlices(t *testing.T) {
	_, err := uncertain.Get([]byte{1, 2, 3}, "wrong")
	be.Nonzero(t, err)

	_, err = uncertain.Get([]int{1, 2, 3}, -1)
	be.Nonzero(t, err)

	_, err = uncertain.Get([]uint{1, 2, 3}, 5)
	be.Nonzero(t, err)

	_, err = uncertain.Get([]uint(nil), 0)
	be.Nonzero(t, err)

	res, err := uncertain.Get([]float32{1, 2, 3}, 1)
	be.NilErr(t, err)
	be.Equal(t, 2, res.(float32))
}

func TestStrings(t *testing.T) {
	_, err := uncertain.Get("string", "segment")
	be.Nonzero(t, err)

	_, err = uncertain.Get("string", float64(100))
	be.Nonzero(t, err)

	_, err = uncertain.Get("string", -1)
	be.Nonzero(t, err)

	_, err = uncertain.Get("string", 20)
	be.Nonzero(t, err)

	res, err := uncertain.Get("string")
	be.NilErr(t, err)
	be.Equal(t, "string", res.(string))

	res, err = uncertain.Get("string", 2)
	be.NilErr(t, err)
	be.Equal(t, 'r', res.(byte))

}

func TestMap(t *testing.T) {
	tm := map[string]interface{}{
		"key1": "val1",
		"key2": nil,
	}

	_, err := uncertain.Get(map[string]string(nil), "key1")
	be.Nonzero(t, err)

	res, err := uncertain.Get(tm, "key1")
	be.NilErr(t, err)
	be.Equal(t, "val1", res.(string))

	res, err = uncertain.Get(tm, "key2")
	be.NilErr(t, err)
	be.Zero(t, res)

	_, err = uncertain.Get(tm, "key3")
	be.Nonzero(t, err)
}

func TestStruct(t *testing.T) {
	ts := testStruct{
		Attr1: "hello world",
		Attr2: 42,
		Attr3: map[string]string{
			"hello": "world",
			"other": "key",
		},
	}

	res, err := uncertain.Get(ts)
	be.NilErr(t, err)
	be.DeepEqual(t, ts, res.(testStruct))

	res, err = uncertain.Get(&ts)
	be.NilErr(t, err)
	be.DeepEqual(t, ts, res.(testStruct))

	res, err = uncertain.Get(&ts, "Attr1")
	be.NilErr(t, err)
	be.Equal(t, "hello world", res.(string))

	res, err = uncertain.Get(ts, "Attr2")
	be.NilErr(t, err)
	be.Equal(t, 42, res.(int))

	res, err = uncertain.Get(ts, "Attr3")
	be.NilErr(t, err)
	be.DeepEqual(t, map[string]string{
		"hello": "world",
		"other": "key",
	}, res.(map[string]string))

	_, err = uncertain.Get(ts, "attr2")
	be.Nonzero(t, err)

	_, err = uncertain.Get(ts, "flarbl")
	be.Nonzero(t, err)

	_, err = uncertain.Get(ts, "flarbl", "hello")
	be.Nonzero(t, err)

	res, err = uncertain.Get(ts, "Attr3", "hello")
	be.NilErr(t, err)
	be.Equal(t, "world", res.(string))

	res, err = uncertain.Get(&ts, "Attr3", "hello")
	be.NilErr(t, err)
	be.Equal(t, "world", res.(string))

	res, err = uncertain.Get(&ts, "Attr3", "other", 0)
	be.NilErr(t, err)
	be.Equal(t, 'k', res.(byte))
}

func TestSliceIteration(t *testing.T) {
	ts := testStruct{
		Attr4: []struct {
			A string
			B string
			C interface{}
		}{
			{A: "First A", B: "First B", C: 123},
			{A: "Second A", B: "Second B", C: "123"},
			{A: "Third A", B: "Third B", C: nil},
			{C: map[string]interface{}{
				"key1": "value1",
				"key2": []string{"one", "two", "three", "four"},
			}},
		},
	}

	res, err := uncertain.Get(ts, "Attr4", 0, "A")
	be.NilErr(t, err)
	be.Equal(t, "First A", res.(string))

	res, err = uncertain.Get(ts, "Attr4", "*", "B")
	be.NilErr(t, err)
	be.DeepEqual(t, []string{"First B", "Second B", "Third B", ""}, res.([]string))

	res, err = uncertain.Get(ts, "Attr4", "*", "C")
	be.NilErr(t, err)
	be.DeepEqual(t, []interface{}{123, "123", nil, map[string]interface{}{
		"key1": "value1",
		"key2": []string{"one", "two", "three", "four"},
	}}, res.([]any))

	res, err = uncertain.Get(ts, "Attr4", "*", "C", "key1")
	be.NilErr(t, err)
	be.DeepEqual(t, []string{"value1"}, res.([]string))
}

func TestNested(t *testing.T) {
	type keyStruct struct{}
	deep := map[interface{}]interface{}{
		"key": "value",
		42:    42,
		keyStruct{}: map[string]interface{}{
			"innerKey": []int{1, 2, 3},
		},
	}

	res, err := uncertain.Get(deep, keyStruct{}, "innerKey", 2)
	be.NilErr(t, err)
	be.Equal(t, 3, res.(int))

	res, err = uncertain.Get(deep, "key", "1")
	be.NilErr(t, err)
	be.Equal(t, 'a', res.(byte))
}

func TestNilPtr(t *testing.T) {
	type nilPtrStruct struct {
		Ohno *string
	}

	nilPtr := nilPtrStruct{Ohno: nil}
	_, err := uncertain.Get(nilPtr, "Ohno", 2)
	be.Nonzero(t, err)
}

func ExampleGet() {
	t := map[string]interface{}{
		"outer": struct{ Inner string }{"value"},
	}
	val, err := uncertain.Get(t, "outer", "Inner")
	fmt.Printf("t[\"outer\"].Inner == \"%s\", err == %s\n", val, err)
}
