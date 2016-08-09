// Package dictionary implements a hash/map/dictionary for educational purposes.
package dictionary

import "container/list"

type (
	// Dictionary is a simple hashed dictionary. It is intended
	// to only store a single type, but does not enforce this explicitly.
	// It is not safe for concurrent use, so users should implement
	// their own locking.
	Dictionary struct {
		numBuckets uint32
		// just use a simple list for our bucket
		// this is not meant for very high performance, just as an example.
		buckets []*list.List
	}

	item struct {
		key   Hasher
		hash  uint32
		value interface{}
	}

	// OptionsFunc is used to set options when creating a new dictionary.
	OptionsFunc func(*Dictionary)

	// EachFunc is the function called on each element when calling Each
	// returning a non-nil error will cause iteration to stop
	EachFunc func(Hasher, interface{}) error

	// Hasher defines interface for keys to be stored in a dictionary.
	Hasher interface {
		// Hash should return a hash of the key. Ideally, this should create
		// a good distribution and avoid collisions.
		Hash() uint32
		// Equal must return true if the receiver is equal to the argument.
		Equal(interface{}) bool
	}
)

// New creates a new dictionary. Options can be set by passing in OptionsFunc
func New(options ...OptionsFunc) *Dictionary {
	d := &Dictionary{
		// 31 is a good choice for a few dozen to a couple hundred keys.
		// We could dynamically resize the number of buckets, but that increases
		// the complexity.
		numBuckets: 31,
	}

	for _, f := range options {
		f(d)
	}

	d.buckets = make([]*list.List, d.numBuckets)
	for i := 0; uint32(i) < d.numBuckets; i++ {
		d.buckets[i] = list.New()
	}
	return d
}

// SetBuckets will set the number of hash buckets.
func SetBuckets(n uint32) func(d *Dictionary) {
	return func(d *Dictionary) {
		d.numBuckets = n
	}
}

func (d *Dictionary) getBucket(key Hasher) (uint32, *list.List) {
	h := key.Hash()
	n := h % d.numBuckets
	return h, d.buckets[n]
}

// Set adds an item to the dictionary. It will replace any existing value.
func (d *Dictionary) Set(key Hasher, val interface{}) {
	h, bucket := d.getBucket(key)

	i := &item{
		hash:  h,
		key:   key,
		value: val,
	}

	// quick exit, bucket is empty
	if bucket.Len() == 0 {
		bucket.PushFront(i)
		return
	}

	for e := bucket.Front(); e != nil; e = e.Next() {
		v := e.Value.(*item)
		// check the hash value first. If these are not equal, then the keys cannot be equal.
		if v.hash == h && key.Equal(v.key) {
			// replace
			e.Value = i
			return
		}
	}

	// key not found, so add it
	bucket.PushFront(i)
}

// helper to get the list and element.
func (d *Dictionary) getElement(key Hasher) (*list.List, *list.Element) {
	h, bucket := d.getBucket(key)
	for e := bucket.Front(); e != nil; e = e.Next() {
		v := e.Value.(*item)
		if v.hash == h && key.Equal(v.key) {
			return bucket, e
		}
	}
	return nil, nil
}

// Get returns an item from the dictionary. The second return value will be
// false if not found.
func (d *Dictionary) Get(key Hasher) (interface{}, bool) {
	bucket, e := d.getElement(key)
	if bucket == nil || e == nil {
		return nil, false
	}
	return e.Value.(*item).value, true

}

// Delete removes an item from the dictionary.  Returns the deleted value.
//
func (d *Dictionary) Delete(key Hasher) (interface{}, bool) {
	bucket, e := d.getElement(key)
	if bucket == nil || e == nil {
		return nil, false
	}
	v := e.Value.(*item).value
	bucket.Remove(e)
	return v, true
}

// Each executes the function on each element. Error returned will be
// any error the EachFunc returned tos top iteration
func (d *Dictionary) Each(f EachFunc) error {
	for _, bucket := range d.buckets {
		for e := bucket.Front(); e != nil; e = e.Next() {
			i := e.Value.(*item)
			if err := f(i.key, i.value); err != nil {
				return err
			}
		}
	}

	return nil
}

// Keys returns all the keys in the hash
func (d *Dictionary) Keys() []Hasher {
	// first calculate the length
	len := 0
	for _, bucket := range d.buckets {
		len = len + bucket.Len()
	}
	keys := make([]Hasher, len)

	i := 0
	for _, bucket := range d.buckets {
		for e := bucket.Front(); e != nil; e = e.Next() {
			keys[i] = e.Value.(*item).key
			i++
		}
	}
	return keys
}
