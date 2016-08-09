package dictionary_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/bakins/dictionary"
	"github.com/stretchr/testify/require"
)

func TestSimpleSet(t *testing.T) {
	d := dictionary.New()
	k := dictionary.StringKey("foo")

	d.Set(k, "bar")
	v, ok := d.Get(k)
	require.NotNil(t, v)
	require.Equal(t, true, ok, "should have found key")
	require.Equal(t, "bar", v.(string), "unexpected value")

	v, ok = d.Get(dictionary.StringKey("bar"))
	require.Nil(t, v)
	require.Equal(t, false, ok, "should not have found key")
}

type entry struct {
	key dictionary.StringKey
	val int
}

func TestSet(t *testing.T) {
	d := dictionary.New()

	entries := make([]entry, 0)
	for i, c := range "abcdefghijklmnopqrstuvwxyz" {
		e := entry{
			key: dictionary.StringKey(c),
			val: i,
		}

		entries = append(entries, e)
		d.Set(e.key, &e)
	}

	for i := range entries {
		j := rand.Intn(i + 1)
		entries[i], entries[j] = entries[j], entries[i]
	}

	for _, e := range entries {
		v, ok := d.Get(e.key)
		require.Equal(t, true, ok, "should have found key")
		require.Equal(t, e.val, v.(*entry).val, "unexpected value")
	}
}

func TestDelete(t *testing.T) {
	d := dictionary.New()
	k := dictionary.StringKey("foo")

	d.Set(k, "bar")
	v, ok := d.Get(k)
	require.NotNil(t, v)
	require.Equal(t, true, ok, "should have found key")
	require.Equal(t, "bar", v.(string), "unexpected value")

	v, ok = d.Delete(dictionary.StringKey("foo"))
	require.NotNil(t, v)
	require.Equal(t, true, ok, "should have found key")
	require.Equal(t, "bar", v.(string), "unexpected value")

	v, ok = d.Delete(dictionary.StringKey("bar"))
	require.Nil(t, v)
	require.Equal(t, false, ok, "should not have found key")
}

type intKey int

func (i intKey) Hash() uint32 {
	if i < 0 {
		i = -i
	}
	if i < math.MaxUint32 {
		return uint32(i)
	}

	// hacky but good enough for a test
	return uint32(i - math.MaxUint32)
}

func (i intKey) Equal(v interface{}) bool {
	return int(i) == int(v.(intKey))
}

func TestSimpleIntSet(t *testing.T) {
	d := dictionary.New()
	k := intKey(99)

	d.Set(k, "bar")
	v, ok := d.Get(k)
	require.NotNil(t, v)
	require.Equal(t, true, ok, "should have found key")
	require.Equal(t, "bar", v.(string), "unexpected value")

	v, ok = d.Get(intKey(1))
	require.Nil(t, v)
	require.Equal(t, false, ok, "should not have found key")
}

type intEntry struct {
	key intKey
	val int
}

func createIntEntries(num int) []intEntry {
	entries := make([]intEntry, num)
	for i := 0; i < num; i++ {
		entries[i] = intEntry{
			key: intKey(i),
			val: i,
		}
	}

	for i := range entries {
		j := rand.Intn(i + 1)
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries
}

func addEntries(d *dictionary.Dictionary, entries []intEntry) {
	for i := 0; i < len(entries); i++ {
		e := &entries[i]
		d.Set(e.key, e)
	}
}

func TestIntSet(t *testing.T) {
	d := dictionary.New()

	entries := createIntEntries(8192)

	addEntries(d, entries)

	for _, e := range entries {
		v, ok := d.Get(e.key)
		require.Equal(t, true, ok, "should have found key")
		require.Equal(t, e.val, v.(*intEntry).val, "unexpected value")
	}
}

func TestEach(t *testing.T) {
	d := dictionary.New()

	keys := []string{"a", "b", "c", "d"}
	entries := make(map[string]string, len(keys))
	for _, k := range keys {
		entries[k] = k
		d.Set(dictionary.StringKey(k), k)
	}

	f := func(h dictionary.Hasher, v interface{}) error {
		k := string(h.(dictionary.StringKey))
		val := v.(string)
		e, ok := entries[k]
		if !ok {
			return fmt.Errorf("did not find %s", k)
		}
		if e != val {
			return fmt.Errorf("bad value - %s - for %s", e, val)
		}
		return nil
	}

	err := d.Each(f)
	require.Nil(t, err)

}

func ExampleNew() {
	d := dictionary.New()
	k := dictionary.StringKey("foo")

	d.Set(k, "bar")
	v, _ := d.Get(k)

	fmt.Println(v.(string))
	// Output: bar
}

func ExampleSetBuckets() {
	d := dictionary.New(dictionary.SetBuckets(997))
	k := dictionary.StringKey("foo")

	d.Set(k, "bar")
	v, _ := d.Get(k)

	fmt.Println(v.(string))
	// Output: bar
}

// TODO: test keys
// TODO: benchmarks of various bucket sizes

func BenchmarkMap(b *testing.B) {
	m := make(map[intKey]int)
	entries := createIntEntries(128)
	for _, e := range entries {
		m[e.key] = e.val
	}
	for n := 0; n < b.N; n++ {
		for _, e := range entries {
			_, _ = m[e.key]
		}
	}
}

func BenchmarkSimple128(b *testing.B) {
	d := dictionary.New()

	entries := createIntEntries(128)

	addEntries(d, entries)

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		for _, e := range entries {
			_, _ = d.Get(e.key)
		}
	}
}

func Benchmark128BucketSize(b *testing.B) {
	d := dictionary.New(dictionary.SetBuckets(127))

	entries := createIntEntries(128)

	addEntries(d, entries)

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		for _, e := range entries {
			_, _ = d.Get(e.key)
		}
	}
}
