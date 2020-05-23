package object

import "hash/fnv"

// HashKey hashes the 3 main types in Xlang
type HashKey struct {
	Type     ObjectType
	Value    uint64
	ValuePtr Object
}

// HashKey Hashing method of a boolean
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// HashKey Hashing method of an integer
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey is a hash mehod of a string
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
