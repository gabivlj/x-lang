package object

import (
	"fmt"
	"strings"
)

// HashPair stores a pair inside a map
type HashPair struct {
	Key   Object
	Value Object
}

// Hashable checks if an object implements the HashKey function
type Hashable interface {
	HashKey() HashKey
}

// HashMap is a hash map in xlang
type HashMap struct {
	Pairs           map[HashKey]HashPair
	UnhashablePairs map[Object]HashPair
}

// Type returns the type of the hash map
func (h *HashMap) Type() ObjectType { return HashObject }

// Inspect inspects the hashmap items
func (h *HashMap) Inspect() string {
	var out strings.Builder

	pairs := make([]string, 0, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	for _, pair := range h.UnhashablePairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString(fmt.Sprintf("{ %s }", strings.Join(pairs, ", ")))

	return out.String()
}
