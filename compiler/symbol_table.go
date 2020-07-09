package compiler

// SymbolScope .
type SymbolScope string

const (
	// LocalScope .
	LocalScope SymbolScope = "LOCAL"
	// GlobalScope .
	GlobalScope SymbolScope = "GLOBAL"
	// BuiltinScope .
	BuiltinScope SymbolScope = "BUILTIN"
	// FreeScope is the scope for variables that are catched by a closure
	FreeScope SymbolScope = "FREE"
)

// Symbol .
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable stores all the symbols
type SymbolTable struct {
	Outer       *SymbolTable
	FreeSymbols []Symbol

	store          map[string]Symbol
	numDefinitions int
}

// NewSymbolTable returns a new table
func NewSymbolTable() *SymbolTable {
	free := []Symbol{}
	return &SymbolTable{store: make(map[string]Symbol), FreeSymbols: free}
}

// Resolve returns a symbol
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	if !ok && s.Outer != nil {
		sym, ok := s.Outer.Resolve(name)
		if !ok {
			return sym, ok
		}
		if sym.Scope == GlobalScope || sym.Scope == BuiltinScope {
			return sym, ok
		}

		return s.defineFree(sym), ok
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

// DefineBuiltin defines a builtin function
func (s *SymbolTable) DefineBuiltin(index int, name string) {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)
	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}
	symbol.Scope = FreeScope
	s.store[original.Name] = symbol
	return symbol
}
