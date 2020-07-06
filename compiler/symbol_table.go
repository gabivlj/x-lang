package compiler

// SymbolScope .
type SymbolScope string

const (
	// GlobalScope .
	GlobalScope SymbolScope = "GLOBAL"
)

// Symbol .
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable stores all the symbols
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

// NewSymbolTable returns a new table
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

// Define returns the new symbol
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
	s.numDefinitions++
	s.store[symbol.Name] = symbol
	return symbol
}

// Resolve returns a symbol
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]

	return symbol, ok
}
