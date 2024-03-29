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

var (
    // operatorSet stores all operator type.
    operatorSet               = map[string]struct{}{token.PLUS: {}, token.MINUS: {}, token.SLASH: {}, token.ASTERISK: {}, token.CIRCUMFLEX: {}}
    correspondingRightBracket = map[string]string{token.LPAREN: token.RPAREN, token.LSQBRACK: token.RSQBRACK, token.LCURBRACK: token.RCURBRACK}
    correspondingLeftBracket  = map[string]string{token.RPAREN: token.LPAREN, token.RSQBRACK: token.LSQBRACK, token.RCURBRACK: token.LCURBRACK}
)

// Parser reads token from the lexer.
type Parser struct {
    root           *ast.Root
    l              *lexer.Lexer
    currToken      *token.Token
    nextToken      *token.Token
    equationCursor int
    result         float64
}

// New creates a new instance of a Parser.
func New(l *lexer.Lexer) *Parser {
    return &Parser{l: l}
}

// Evaluate takes input and calculates the result.
func (p *Parser) Evaluate(input string) (float64, error) {
    p.input(input)
    err := p.parsePrompt()
    if err != nil {
        return 0, err
    }
    n, err := p.parseEquation(0)
    if err != nil {
        return 0, err
    }
    result, err := ast.Evaluate(n)
    p.result = result
    return result, err
}

// ClearPreviousAns resets the stored result to 0.
func (p *Parser) ClearPreviousAns() {
    p.result = 0
}

// input takes input data and send it to the lexer.
// It frees the root and set the equationCursor to 0 before taking input.
func (p *Parser) input(data string) {
    p.root = nil
    p.equationCursor = 0
    p.l.Input(data)
}

// peekToken gets a copy of the next token.
func (p *Parser) peekToken() token.Token {
    return *p.nextToken
}

// readNextToken returns the nextToken read from the parser.
func (p *Parser) readNextToken() *token.Token {
    return p.l.ReadNextToken()
}

// parsePrompt evaluates calculator prompt and tell whether it's valid.
// If the given input is valid, the equation node is then stored for further evaluation.
// If it isn't, parsePrompt returns an error.
func (p *Parser) parsePrompt() error {
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
    if !hasBalanceBrackets(p.root.EquationTokens) {
        return ErrEquation
    }

    return nil
}

// peekEquationToken retrieves the next token from equationTokens without advancing the cursor.
// If the cursor is at the end of the token list, it returns nil to indicate that there are no more tokens.
func (p *Parser) peekEquationToken() *token.Token {
    if p.equationCursor >= len(p.root.EquationTokens) {
        return nil
    }
    return p.root.EquationTokens[p.equationCursor]
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

// parseEquation parses the equation stored in the Parser into an *ast.Node.
func (p *Parser) parseEquation(minbp int) (*ast.Node, error) {
    // Left hand side. Calling nextEquationToken returns the currToken the pointer is pointing at then advances it.
    lhsTok := p.nextEquationToken()

    var lhs *ast.Node
    switch {
    case isInt(lhsTok):
        lhsVal, _ := strconv.Atoi(lhsTok.Literal)
        lhs = ast.New(lhsTok, float64(lhsVal), true, "", false, nil, nil)
    case isFloat(lhsTok):
        lhsVal, _ := strconv.ParseFloat(lhsTok.Literal, 32)
        lhs = ast.New(lhsTok, lhsVal, true, "", false, nil, nil)

    case isOperator(lhsTok):
        rbp := prefixBindingPower(lhsTok)
        // Scenario: Unknown operator, we shouldn't have operators other than '+' and '-'.
        if rbp == 0 {
            return nil, ErrEquation
        }

        // Scenario: Something wrong happened when parsing a deeper node, like '5 + 6 *'.
        rightChild, err := p.parseEquation(rbp)
        if err != nil {
            return nil, err
        }

        lhs = ast.New(lhsTok, 0, false, lhsTok.Literal, true, nil, rightChild)

    case isLeftBracket(lhsTok):
        var err error

        // We know that an equation expression should exist within the bracket (no matter valid or not).
        lhs, err = p.parseEquation(0)
        if err != nil {
            // Cases like [(12 + 3 * 6] + 1 might happen.
            // We need to inherit the error from the last recursive call and break.
            return nil, err
        }

        if p.peekEquationToken() == nil {
            return nil, ErrEquation
        } else if p.peekEquationToken().LexicalType != correspondingRightBracket[lhsTok.LexicalType] {
            // Closing brackets doesn't match.
            return nil, ErrEquation
        } else {
            // Consume the correct right parenthesis.
            p.nextEquationToken()
        }
    case isAns(lhsTok):
        lhs = ast.New(lhsTok, p.result, true, "", false, nil, nil)

    default:
        // Scenario: Missing an integer, like '' or '5 + '.
        // Or if the token we encounter is an unknown type.
        return nil, ErrEquation
    }

    for {
        op := p.peekEquationToken()
        if op == nil {
            // If there is no more operators in equation, break out of the loop.
            break
        }

        if isRightBracket(op) {
            // If it's a right bracket, break out of the loop.
            break
        }

        // Scenario: Missing operator between integer tokens, like '5 25'.
        // This also deals with something like '0)'.
        if !isOperator(op) {
            return nil, ErrEquation
        }

        // Get the binding power for the current operator.
        lbp, rbp := infixBindingPower(op)
        // Scenario: Unknown operator.
        if lbp == 0 || rbp == 0 {
            return nil, ErrEquation
        }

        if lbp < minbp {
            // When left binding power is smaller than the minimum binding power
            // we have higher precedence before than the one we currently encounter, break.
            break
        }

        // Advance the equation cursor, this results in the cursor pointing at the operator.
        p.nextEquationToken()

        // Fetch the right hand side node recursively.
        rhs, err := p.parseEquation(rbp)
        if err != nil {
            return nil, err
        }

        // Updates lhs for the next iteration.
        lhs = formEquation(op, lhs, rhs)
    }
    return lhs, nil
}

// formEquation creates a new ast.Node representing an operator and its operands.
func formEquation(op *token.Token, lhs *ast.Node, rhs *ast.Node) *ast.Node {
    operatorNode := ast.New(op, 0, false, op.Literal, true, lhs, rhs)
    return operatorNode
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
    case token.CIRCUMFLEX:
        return 6, 7
    }
    return 0, 0
}

// prefixBindingPower returns the right binding powers for prefix operators like '+' and '-'.
func prefixBindingPower(operatorToken *token.Token) int {
    switch operatorToken.LexicalType {
    case token.PLUS:
        return 5
    case token.MINUS:
        return 5
    default:
        return 0
    }
}

// isInt checks whether a token is an integer.
func isInt(tok *token.Token) bool {
    if tok != nil {
        return tok.LexicalType == token.INT
    }
    return false
}

// isFloat checks whether a token is a float.
func isFloat(tok *token.Token) bool {
    if tok != nil {
        return tok.LexicalType == token.FLOAT
    }
    return false
}

// isOperator checks whether a token is an operator.
func isOperator(tok *token.Token) bool {
    if tok != nil {
        _, ok := operatorSet[tok.LexicalType]
        return ok
    }
    return false
}

// isLeftBracket checks whether a token is a left bracket.
func isLeftBracket(tok *token.Token) bool {
    if tok == nil {
        return false
    }

    _, ok := correspondingRightBracket[tok.LexicalType]
    return ok
}

// isRightBracket checks whether a token is a right bracket.
func isRightBracket(tok *token.Token) bool {
    if tok == nil {
        return false
    }
    _, ok := correspondingLeftBracket[tok.LexicalType]
    return ok
}

// isAns checks whether a token is an 'ans' token.
func isAns(tok *token.Token) bool {
    if tok != nil {
        return tok.LexicalType == token.ANS
    }
    return false
}

// hasBalanceBrackets checks whether the tokens has balances parentheses.
func hasBalanceBrackets(tokens []*token.Token) bool {
    stack := make([]*token.Token, 0)
    pop := func(stack *[]*token.Token) {
        if len(*stack) != 0 {
            *stack = (*stack)[0 : len(*stack)-1]
        }
    }

    // Check if the opening and closing parentheses matches.
    // Check invalid brackets like [(]).
    // Check if there are adjacent opening and closing brackets, (), [] and {} are invalid.
    for n, tok := range tokens {
        switch {
        case isLeftBracket(tok):
            stack = append(stack, tok)
        case isRightBracket(tok):
            // Check if the last token is the current tokens' corresponding left bracket, i.e. ()23 + 5 or 23 + 5().
            if n > 0 && tokens[n-1].LexicalType != correspondingLeftBracket[tok.LexicalType] {
                // If the stack length is 0, we have a redundant closing bracket.
                // If the type of last token in the stack isn't the corresponding left bracket, we have invalid bracket grammar, i.e. [(]).
                if len(stack) == 0 || stack[len(stack)-1].LexicalType != correspondingLeftBracket[tok.LexicalType] {
                    return false
                }
                pop(&stack)
            }
        }
    }

    // If the length of the stack isn't 0, we have opening brackets that aren't closed.
    return len(stack) == 0
}
