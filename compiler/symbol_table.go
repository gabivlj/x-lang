package compiler

// SymbolScope .
type SymbolScope string

const (
	// LocalScope .
	LocalScope SymbolScope = "LOCAL"
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
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int
}

// NewSymbolTable returns a new table
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

// // Define returns the new symbol
// func (s *SymbolTable) Define(name string) Symbol {
// 	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
// 	s.numDefinitions++
// 	s.store[symbol.Name] = symbol
// 	return symbol
// }

// Resolve returns a symbol
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	if !ok && s.Outer != nil {
		return s.Outer.Resolve(name)
	}
	return symbol, ok
}

// NewEnclosedSymbolTable returns a new symbol table
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

// Define a new symbol
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}
