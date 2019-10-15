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
}

func TestGet(t *testing.T) {
	res, err := uncertain.Get(func() {}, "segment")
	require.Error(t, err)

	res, err = uncertain.Get(make(chan string), "segment")
	require.Error(t, err)

	res, err = uncertain.Get(nil)
	require.NoError(t, err)
	require.Nil(t, res, "Get(nil) == nil")

	res, err = uncertain.Get(nil, "segment")
	require.Error(t, err)

	res, err = uncertain.Get(nil, 1)
	require.Error(t, err)

	res, err = uncertain.Get([]byte{1, 2, 3}, "wrong")
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

	res, err = uncertain.Get("string", "segment")
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

	ts := testStruct{
		Attr1: "hello world",
		Attr2: 42,
		Attr3: map[string]string{
			"hello": "world",
			"other": "key",
		},
	}

	res, err = uncertain.Get(ts)
	require.NoError(t, err)
	require.Exactly(t, ts, res, "Must be the test struct")

	res, err = uncertain.Get(&ts)
	require.NoError(t, err)
	require.Exactly(t, &ts, res, "Must be the test struct ptr")

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

	type keyStruct struct{}
	deep := map[interface{}]interface{}{
		"key": "value",
		42:    42,
		keyStruct{}: map[string]interface{}{
			"innerKey": []int{1, 2, 3},
		},
	}

	res, err = uncertain.Get(deep, keyStruct{}, "innerKey", 2)
	require.NoError(t, err)
	require.Exactly(t, 3, res, "Must be 3")

	res, err = uncertain.Get(deep, "key", "1")
	require.NoError(t, err)
	require.Exactly(t, byte('a'), res, "Must be 'a'")
}

func ExampleGet() {
	t := map[string]interface{}{
		"outer": struct{ Inner string }{"value"},
	}
	val, err := uncertain.Get(t, "outer", "Inner")
	fmt.Printf("t[\"outer\"].Inner == \"%s\", err == %s\n", val, err)
}
