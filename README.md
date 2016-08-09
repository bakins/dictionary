# dictionary
[![Build Status](https://travis-ci.org/bakins/dictionary.svg?branch=master)](https://travis-ci.org/bakins/dictionary)
[![GoDoc](https://godoc.org/github.com/bakins/dictionary?status.png)](https://godoc.org/github.com/bakins/dictionary)

Simple
[dictionary/hash-table](https://en.wikipedia.org/wiki/Hash_table) in
Go for education/testing.  It uses an array of double-linked list for
the actual storage.  This is a good compromie between performance,
memory usage, and complexity.  The number of buckets can be set at
creation time.

The [tests](./dictionary_test.go) provide examples of usage.



