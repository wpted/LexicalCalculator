package parser

import (
    "LexicalCalculator/ast"
    "LexicalCalculator/lexer"
    "LexicalCalculator/token"
    "errors"
    "strconv"
)

var (
    ErrPrompt       = errors.New("error incorrect prompt")
    ErrOpeningQuote = errors.New("error missing prompt opening quote")
    ErrClosingQuote = errors.New("error missing prompt closing quote")
    ErrEquation     = errors.New("error equation format")
)

// Parser reads token from the lexer.
type Parser struct {
    root           *ast.Root
    l              *lexer.Lexer
    currToken      *token.Token
    nextToken      *token.Token
    equationCursor int
}

// New creates a new instance of a Parser.
func New(l *lexer.Lexer) *Parser {
    return &Parser{l: l}
}

// Input takes input data and send it to the lexer.
// It frees the root and set the equationCursor to 0 before taking input.
func (p *Parser) Input(data string) {
    p.root = nil
    p.equationCursor = 0
    p.l.Input(data)
}

// Parse evaluates calculator prompt and tell whether it's valid.
// If the given input is valid, the equation node is then stored for further evaluation.
// If it isn't, Parse returns an empty root and an error.
func (p *Parser) Parse() error {
    p.root = new(ast.Root)
    p.root.EquationTokens = make([]*token.Token, 0)
    // Read the first token and check if it's started with 'calc'.
    firstToken := p.readNextToken()
    if firstToken.LexicalType == token.EOF {
        // This happens when the prompt is empty.
        return ErrPrompt
    }
    if firstToken.Literal != "calc" && firstToken.LexicalType != token.CALC {
        return ErrPrompt
    }
    // After checking, assign the first token to root.
    p.root.Token = firstToken

    // Assign it to the currToken.
    p.currToken = firstToken

    // Read the next token. The following token should be a single quote token.
    nextToken := p.readNextToken()
    if nextToken.LexicalType == token.EOF {
        // This happens when there is no following string after calc.
        return ErrPrompt
    }
    if nextToken.Literal != "'" && nextToken.LexicalType != token.SINGLEQUOTE {
        return ErrOpeningQuote
    }

    // After checking, assign it to the nextToken.
    p.nextToken = nextToken

    for {
        if p.nextToken.LexicalType == token.EOF {
            if p.currToken.LexicalType != token.SINGLEQUOTE {
                return ErrClosingQuote
            }
            break
        }

        // If the next token is not token.EOF, advance the currToken.
        p.currToken = p.nextToken

        newToken := p.readNextToken()
        p.nextToken = newToken

        // Store equation nodes.
        if p.currToken.LexicalType != token.SINGLEQUOTE {
            p.root.EquationTokens = append(p.root.EquationTokens, p.currToken)
        }
    }

    return nil
}

// peekToken gets a copy of the next token.
func (p *Parser) peekToken() token.Token {
    return *p.nextToken
}

// readNextToken returns the nextToken read from the parser.
func (p *Parser) readNextToken() *token.Token {
    return p.l.ReadNextToken()
}

// nextEquationToken retrieves the next token from equationTokens list and advances the cursor.
// If the cursor is at the end of the token list, it returns nil to indicate that there are no more tokens.
func (p *Parser) nextEquationToken() *token.Token {
    if p.equationCursor >= len(p.root.EquationTokens) {
        return nil
    }
    tok := p.root.EquationTokens[p.equationCursor]
    p.equationCursor++
    return tok
}

// peekEquationToken retrieves the next token from equationTokens without advancing the cursor.
// If the cursor is at the end of the token list, it returns nil to indicate that there are no more tokens.
func (p *Parser) peekEquationToken() *token.Token {
    if p.equationCursor >= len(p.root.EquationTokens) {
        return nil
    }
    return p.root.EquationTokens[p.equationCursor]
}

// parseEquation parses the equation stored in the Parser into an *ast.Node.
func (p *Parser) parseEquation(minbp int) *ast.Node {
    // Left hand side. Calling nextEquationToken returns the currToken the pointer is pointing at then advances it.
    lhsTok := p.nextEquationToken()

    lhsVal, _ := strconv.Atoi(lhsTok.Literal)
    lhs := ast.New(lhsTok, lhsVal, true, "", false, nil, nil)

    for {
        op := p.peekEquationToken()
        if op == nil {
            // If there is no more operators ( or following tokens ), break out of the loop.
            break
        }
        // Get the binding power for the current operator.
        lbp, rbp := infixBindingPower(op)
        if lbp < minbp {
            // When left binding power is smaller than the minimum binding power
            // we have higher precedence before than the one we currently encounter, break.
            break
        }

        // Advance the equation cursor, this results in the cursor pointing at the operator.
        p.nextEquationToken()

        // Fetch the right hand side node recursively.
        rhs := p.parseEquation(rbp)

        // Updates lhs for the next iteration.
        lhs = formBinaryTree(op, lhs, rhs)
    }

    return lhs
}

// infixBindingPower returns the left and right binding powers for different operators.
func infixBindingPower(operatorToken *token.Token) (int, int) {
    switch operatorToken.LexicalType {
    case token.PLUS:
        return 1, 2
    case token.MINUS:
        return 1, 2
    case token.ASTERISK:
        return 3, 4
    case token.SLASH:
        return 3, 4
    }
    return 0, 0
}

// formBinaryTree creates a new ast.Node representing an operator and its operands.
func formBinaryTree(op *token.Token, lhs *ast.Node, rhs *ast.Node) *ast.Node {
    operatorNode := ast.New(op, 0, false, op.Literal, true, lhs, rhs)
    return operatorNode
}
