package token

const (
    CALC = "CALC"

    LEFTQUOTE  = "LEFTQUOTE"
    RIGHTQUOTE = "RIGHTQUOTE"

    INT = "INT"

    PLUS     = "+"
    MINUS    = "-"
    ASTERISK = "*"
    SLASH    = "/"
)

// Token is the result of after parsing input with a lexer.
type Token struct {
    Literal     string
    LexicalType string
}

// New creates a new Token.
func New(lexicalType string, tokenLiteral string) *Token {
    return &Token{
        LexicalType: lexicalType,
        Literal:     tokenLiteral,
    }
}
