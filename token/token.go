package token

const (
    CALC = "CALC"
    ANS  = "ANS"

    SINGLEQUOTE = "'"
    LPAREN      = "("
    RPAREN      = ")"
    LSQBRACK    = "["
    RSQBRACK    = "]"
    LCURBRACK   = "{"
    RCURBRACK   = "}"

    INT   = "INT"
    FLOAT = "FLOAT"

    PLUS       = "+"
    MINUS      = "-"
    ASTERISK   = "*"
    SLASH      = "/"
    CIRCUMFLEX = "^"

    UNKNOWN = "UNKNOWN"
    EOF     = "EOF"
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
