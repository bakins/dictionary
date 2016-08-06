// Package dictionary implements a hash/map/dictionary for educational purposes.
package dictionary

import (
	"container/list"
	"sync"
)

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
		// we could keep the hash value as an optimization
		// if compare function was considered to be expensive
		key   Hasher
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
		// a good distribution
		Hash() uint32
		// Compare must return -1 if the item is less than the argument. 0 if they are equal
		// and 1 if the argument is greater than the item.
		Compare(interface{}) int
	}
)

var itemPool = sync.Pool{
	New: func() interface{} {
		return &item{}
	},
}

// New creates a new dictionary. Options can be set by passing in OptionsFunc
func New(options ...OptionsFunc) *Dictionary {
	d := &Dictionary{
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

func (d *Dictionary) getBucket(key Hasher) *list.List {
	n := key.Hash() % d.numBuckets
	return d.buckets[n]
}

// Set adds an item to the dictionary. It will replace any existing value.
func (d *Dictionary) Set(key Hasher, val interface{}) {
	bucket := d.getBucket(key)

	i := itemPool.Get().(*item)
	i.key = key
	i.value = val

	// quick exit, bucket is empty
	if bucket.Len() == 0 {
		bucket.PushFront(i)
		return
	}

	for e := bucket.Front(); e != nil; e = e.Next() {
		v := e.Value.(*item)
		// we use compare so that we do not have to traverse the entire list.
		// this is a very minor optimization.
		switch key.Compare(v.key) {
		case -1:
			// inserted key is "less than" the element, so insert it before.
			// this allows the list to stay sorted
			bucket.InsertBefore(i, e)
			return
		case 0:
			// replace
			itemPool.Put(v)
			e.Value = i
			return
		}
	}

	// if we make it to here, then insert at end as the key is greater than anything else
	bucket.PushBack(i)
}

// helper to get the list and element.
func (d *Dictionary) getElement(key Hasher) (*list.List, *list.Element) {
	bucket := d.getBucket(key)
	for e := bucket.Front(); e != nil; e = e.Next() {
		item := e.Value.(*item)
		switch key.Compare(item.key) {
		case 0:
			return bucket, e
		case -1:
			// this key is not in the list
			return nil, nil
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
	i := e.Value.(*item)
	v := i.value
	itemPool.Put(i)
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
