package dictionary

import (
	"hash/crc32"
	"strings"
)

// StringKey is a convinience type for using strings as keys in a dictionary
type StringKey string

func (s StringKey) Hash() uint32 {
	return crc32.ChecksumIEEE([]byte(string(s)))
}

func (s StringKey) Compare(v interface{}) int {
	return strings.Compare(string(s), string(v.(StringKey)))
}

func (s StringKey) String() string {
	return string(s)
}
