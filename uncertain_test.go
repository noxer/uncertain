package uncertain_test

import (
	"fmt"
	"testing"

	"github.com/noxer/uncertain"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Attr1 string
	Attr2 int
	Attr3 map[string]string
	Attr4 interface{}
}

func TestInvalid(t *testing.T) {
	_, err := uncertain.Get(func() {}, "segment")
	require.Error(t, err)

	_, err = uncertain.Get(make(chan string), "segment")
	require.Error(t, err)
}

func TestNil(t *testing.T) {
	res, err := uncertain.Get(nil)
	require.NoError(t, err)
	require.Nil(t, res, "Get(nil) == nil")

	res, err = uncertain.Get(nil, "segment")
	require.Error(t, err)

	res, err = uncertain.Get(nil, 1)
	require.Error(t, err)
}

func TestSlices(t *testing.T) {
	res, err := uncertain.Get([]byte{1, 2, 3}, "wrong")
	require.Error(t, err)

	res, err = uncertain.Get([]int{1, 2, 3}, -1)
	require.Error(t, err)

	res, err = uncertain.Get([]uint{1, 2, 3}, 5)
	require.Error(t, err)

	res, err = uncertain.Get([]uint(nil), 0)
	require.Error(t, err)

	res, err = uncertain.Get([]float32{1, 2, 3}, 1)
	require.NoError(t, err)
	require.Exactly(t, float32(2), res, "Must be 2.0f")
}

func TestStrings(t *testing.T) {
	res, err := uncertain.Get("string", "segment")
	require.Error(t, err)

	res, err = uncertain.Get("string", float64(100))
	require.Error(t, err)

	res, err = uncertain.Get("string", -1)
	require.Error(t, err)

	res, err = uncertain.Get("string", 20)
	require.Error(t, err)

	res, err = uncertain.Get("string")
	require.NoError(t, err)
	require.Exactly(t, "string", res, "Must be 'string'")

	res, err = uncertain.Get("string", 2)
	require.NoError(t, err)
	require.Exactly(t, byte('r'), res, "Must be 'r'")

}

func TestMap(t *testing.T) {
	tm := map[string]interface{}{
		"key1": "val1",
		"key2": nil,
	}

	res, err := uncertain.Get(map[string]string(nil), "key1")
	require.Error(t, err)

	res, err = uncertain.Get(tm, "key1")
	require.NoError(t, err)
	require.Exactly(t, "val1", res, "Must be 'val1'")

	res, err = uncertain.Get(tm, "key2")
	require.NoError(t, err)
	require.Exactly(t, nil, res, "Must be nil")

	res, err = uncertain.Get(tm, "key3")
	require.Error(t, err)
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
	require.NoError(t, err)
	require.Exactly(t, ts, res, "Must be the test struct")

	res, err = uncertain.Get(&ts)
	require.NoError(t, err)
	require.Exactly(t, ts, res, "Must be the dereferenced test struct")

	res, err = uncertain.Get(&ts, "Attr1")
	require.NoError(t, err)
	require.Exactly(t, "hello world", res, "Must be 'hello world'")

	res, err = uncertain.Get(ts, "Attr2")
	require.NoError(t, err)
	require.Exactly(t, 42, res, "Must be 42")

	res, err = uncertain.Get(ts, "Attr3")
	require.NoError(t, err)
	require.Exactly(t, map[string]string{
		"hello": "world",
		"other": "key",
	}, res, "Must be the map")

	res, err = uncertain.Get(ts, "attr2")
	require.Error(t, err)

	res, err = uncertain.Get(ts, "flarbl")
	require.Error(t, err)

	res, err = uncertain.Get(ts, "flarbl", "hello")
	require.Error(t, err)

	res, err = uncertain.Get(ts, "Attr3", "hello")
	require.NoError(t, err)
	require.Exactly(t, "world", res, "Must be 'world'")

	res, err = uncertain.Get(&ts, "Attr3", "hello")
	require.NoError(t, err)
	require.Exactly(t, "world", res, "Must be 'world'")

	res, err = uncertain.Get(&ts, "Attr3", "other", 0)
	require.NoError(t, err)
	require.Exactly(t, byte('k'), res, "Must be 'k'")
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
	require.NoError(t, err)
	require.Exactly(t, "First A", res, "Must be 'First A'")

	res, err = uncertain.Get(ts, "Attr4", "*", "B")
	require.NoError(t, err)
	require.Exactly(t, []string{"First B", "Second B", "Third B", ""}, res, "Must be '[]string{\"First B\", \"Second B\", \"Third B\", \"\"}'")

	res, err = uncertain.Get(ts, "Attr4", "*", "C")
	require.NoError(t, err)
	require.Exactly(t, []interface{}{123, "123", nil, map[string]interface{}{
		"key1": "value1",
		"key2": []string{"one", "two", "three", "four"},
	}}, res, "Must be '[]interface{}{123, \"123\", nil, map[string]interface{}{...}}'")

	res, err = uncertain.Get(ts, "Attr4", "*", "C", "key1")
	require.NoError(t, err)
	require.Exactly(t, []string{"value1"}, res, "Must be '[]string{\"value1\"}'")
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
	require.NoError(t, err)
	require.Exactly(t, 3, res, "Must be 3")

	res, err = uncertain.Get(deep, "key", "1")
	require.NoError(t, err)
	require.Exactly(t, byte('a'), res, "Must be 'a'")
}

func TestNilPtr(t *testing.T) {
	type nilPtrStruct struct {
		Ohno *string
	}

	nilPtr := nilPtrStruct{Ohno: nil}
	_, err := uncertain.Get(nilPtr, "Ohno", 2)
	require.Error(t, err)
}

func ExampleGet() {
	t := map[string]interface{}{
		"outer": struct{ Inner string }{"value"},
	}
	val, err := uncertain.Get(t, "outer", "Inner")
	fmt.Printf("t[\"outer\"].Inner == \"%s\", err == %s\n", val, err)
}
