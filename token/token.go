package token

// Type is the token type itself, not the best for performance reasons, but for this will do the job. We define a token as somethign that we should take into account when parsing the code.
type TypeToken string

// Token is the token that we are gonna parse
type Token struct {
	Type    TypeToken
	Literal string
}

const (
	// ILLEGAL

	ILLEGAL = TypeToken("ILLEGAL") // Unknown token
	EOF     = TypeToken("EOF")     // End of file

	// Identifier + Literals

	IDENT = TypeToken("IDENT") // VARIABLE NAME
	INT   = TypeToken("INT")   // 12345

	// Operators

	COMMA    = TypeToken(",")
	PLUS     = TypeToken("+")
	MINUS    = TypeToken("-")
	BANG     = TypeToken("!")
	ASTERISK = TypeToken("*")
	SLASH    = TypeToken("/")
	LT       = TypeToken("<")
	GT       = TypeToken(">")
	EQ       = TypeToken("==")
	NOTEQ    = TypeToken("!=")

	// Delimiters

	LPAREN = TypeToken("(")
	RPAREN = TypeToken(")")
	LBRACE = TypeToken("{")
	RBRACE = TypeToken("}")

	SEMICOLON = TypeToken(";")
	COLON     = TypeToken(",")
	ASSIGN    = TypeToken("=")

	// Keywords

	FUNCTION = TypeToken("FUNCTION")
	LET      = TypeToken("LET")
	TRUE     = TypeToken("TRUE")
	FALSE    = TypeToken("FALSE")
	IF       = TypeToken("IF")
	ELSE     = TypeToken("ELSE")
	RETURN   = TypeToken("RETURN")

	STRING = TypeToken("STRING")
)

var keywords = map[string]TypeToken{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent Looks up in the keywords table if its a keyword, if its not it will return IDENT as a TypeToken
func LookupIdent(ident string) TypeToken {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
