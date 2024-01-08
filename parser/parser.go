package parser

import (
    "LexicalCalculator/ast"
    "LexicalCalculator/lexer"
    "LexicalCalculator/token"
    "errors"
)

var (
    ErrPrompt       = errors.New("error incorrect prompt")
    ErrOpeningQuote = errors.New("error missing prompt opening quote")
    ErrClosingQuote = errors.New("error missing prompt closing quote")
)

// Parser reads token from the lexer.
type Parser struct {
    l         *lexer.Lexer
    currToken *token.Token
    nextToken *token.Token
}

// New creates a new instance of a Parser.
func New(l *lexer.Lexer) *Parser {
    return &Parser{l: l}
}

// Input takes input data and send it to the lexer.
func (p *Parser) Input(data string) {
    p.l.Input(data)
}

// Parse evaluates calculator prompt and tell whether it's valid.
// If the given input is valid, the equation node is then stored for further evaluation.
// If it isn't, Parse returns an empty root and an error.
func (p *Parser) Parse() (*ast.Root, error) {
    root := new(ast.Root)
    // Read the first token and check if it's started with 'calc'.
    firstToken := p.l.ReadNextToken()
    if firstToken.LexicalType == token.EOF {
        // This happens when the prompt is empty.
        return nil, ErrPrompt
    }
    if firstToken.Literal != "calc" && firstToken.LexicalType != token.CALC {
        return nil, ErrPrompt
    }
    // After checking, assign the first token to root.
    root.Token = firstToken

    // Assign it to the currToken.
    p.currToken = firstToken

    // Read the next token. The following token should be a single quote token.
    nextToken := p.l.ReadNextToken()
    if nextToken.LexicalType == token.EOF {
        // This happens when there is no following string after calc.
        return nil, ErrPrompt
    }
    if nextToken.Literal != "'" && nextToken.LexicalType != token.SINGLEQUOTE {
        return nil, ErrOpeningQuote
    }

    // After checking, assign it to the nextToken.
    p.nextToken = nextToken

    for {
        if p.nextToken.LexicalType == token.EOF {
            if p.currToken.LexicalType != token.SINGLEQUOTE {
                return nil, ErrClosingQuote
            }
            break
        }
        // If the next token is not token.EOF, advance the currToken.
        p.currToken = p.nextToken

        newToken := p.l.ReadNextToken()
        p.nextToken = newToken

        // TODO: Form equation node.
    }

    return root, nil
}

func evaluate(equation string) int {
    // Scenarios
    // 1. 6 ( One single integer )
    // 2. 6 + 5 ( Binary Expression )

    return 0
}
