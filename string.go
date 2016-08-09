package dictionary

import "hash/crc32"

// StringKey is a convinience type for using strings as keys in a dictionary
type StringKey string

// Hash generates a hash for the string using crc32
func (s StringKey) Hash() uint32 {
	return crc32.ChecksumIEEE([]byte(string(s)))
}

// Compare uses the stdlib strings.Compare to compare two string keys
func (s StringKey) Equal(v interface{}) bool {
	return string(s) == string(v.(StringKey))
}

// String returns the string value of the key
func (s StringKey) String() string {
	return string(s)
}
