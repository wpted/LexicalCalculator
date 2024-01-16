/*
Package lexer implements utilities of a lexical parser.
Lexer takes input and parse input into token ( if the lexer understands the input ).
Should be able to reuse, not creating a lexer everytime we receive an input.
*/
package lexer

import (
    "LexicalCalculator/token"
    "bytes"
    "errors"
    "fmt"
)

var (
    ErrInvalidLiteral = errors.New("error token not valid")
)

// Lexer is the lexical analyzer used in the calculator.
// It's supposed to be called by a parser and will lazily return the next token on the fly.
// Lexer resets inputBuffer and resets field currPosition and nextPosition after receiving a new input.
type Lexer struct {
    inputBuffer  *bytes.Buffer
    bufferLength int
    currPosition int
    nextPosition int
    errors       []error
}

// New creates a new Lexer.
func New() *Lexer {
    return &Lexer{
        inputBuffer:  new(bytes.Buffer),
        bufferLength: 0,
        currPosition: -1,
        nextPosition: 0,
    }
}

// Input writes data into the lexer. If the buffer is not empty, it frees the buffer first.
func (l *Lexer) Input(data string) {
    l.free()
    // Write appends the data to the buffer, growing the buffer as needed. The return value n is the length of p; err is always nil.
    // If the buffer becomes too large, Write will panic with ErrTooLarge.
    writtenLength, _ := l.inputBuffer.Write([]byte(data))
    l.bufferLength = writtenLength
}

// free resets the buffer to be empty, but it retains the underlying storage for use by future writes.
func (l *Lexer) free() {
    l.inputBuffer.Reset()
    l.bufferLength = 0
    l.currPosition = -1
    l.nextPosition = 0
}

// ReadNextToken returns a token after every read, skipping all encountered white spaces.
func (l *Lexer) ReadNextToken() *token.Token {
    tok := token.New(token.EOF, token.EOF)

    if l.nextPosition > l.bufferLength {
        return tok
    }

    next := l.inputBuffer.Next(1)
    l.currPosition++
    l.nextPosition++

    // When encounter a white space, advance the pointer.
    if len(next) > 0 {
        ok := isWhiteSpace(next[0])
        for ok {
            next = l.inputBuffer.Next(1)
            l.currPosition++
            l.nextPosition++
            ok = isWhiteSpace(next[0])
        }
    }

    switch string(next) {
    case "":
        return tok
    case "'":
        tok = token.New(token.SINGLEQUOTE, string(next))
    case "(":
        tok = token.New(token.LPAREN, string(next))
    case ")":
        tok = token.New(token.RPAREN, string(next))
    case "[":
        tok = token.New(token.LSQBRACK, string(next))
    case "]":
        tok = token.New(token.RSQBRACK, string(next))
    case "{":
        tok = token.New(token.LCURBRACK, string(next))
    case "}":
        tok = token.New(token.RCURBRACK, string(next))
    case "+":
        tok = token.New(token.PLUS, string(next))
    case "-":
        tok = token.New(token.MINUS, string(next))
    case "*":
        tok = token.New(token.ASTERISK, string(next))
    case "/":
        tok = token.New(token.SLASH, string(next))
    default:
        // We handle integers and 'calc' here.
        if isDigit(next[0]) {
            // The current ch is next[0].
            // What we do here is to advance the pointer and fetch the entire token.
            literal, isInt, err := l.readNumber(next[0])
            if err != nil {
                tok = token.New(token.UNKNOWN, literal)
            } else {
                if isInt {
                    tok = token.New(token.INT, literal)
                } else {
                    tok = token.New(token.FLOAT, literal)
                }
            }
        } else if isLetter(next[0]) {
            literal := l.readIdentifier(next[0])
            switch literal {
            case "calc":
                tok = token.New(token.CALC, literal)
            case "ans":
                tok = token.New(token.ANS, literal)
            default:
                tok = token.New(token.UNKNOWN, literal)
            }
        } else {
            // unknown, append the lexer error.
            tok = token.New(token.UNKNOWN, string(next))
        }
    }
    return tok
}

// readNumber returns the integer literal of a number.
// It advances the pointer until it's at the end of an input or when the next character isn't a digit.
func (l *Lexer) readNumber(first byte) (literal string, isInt bool, err error) {
    decimalSeparatorCount := 0

    numberLiteral := make([]byte, 0)
    numberLiteral = append(numberLiteral, first)
    for {
        peekedToken, _ := peekBuffer(*l.inputBuffer, 1)
        if isDigit(peekedToken) || isDecimalSeparator(peekedToken) {
            next := l.inputBuffer.Next(1)
            l.currPosition++
            l.nextPosition++
            numberLiteral = append(numberLiteral, next[0])

            if isDecimalSeparator(peekedToken) {
                decimalSeparatorCount++
                peekedToken, _ = peekBuffer(*l.inputBuffer, 1)
                if !isDigit(peekedToken) {
                    // When there's no digit after the decimal separator.
                    err = ErrInvalidLiteral
                }
            }
        } else {
            break
        }
    }
    switch {
    case err != nil, decimalSeparatorCount > 1:
        return string(numberLiteral), false, ErrInvalidLiteral
    case decimalSeparatorCount == 1:
        return string(numberLiteral), false, nil
    default:
        return string(numberLiteral), true, nil
    }
}

// readNumber returns the literal of an identifier.
// It advances the pointer until it's at the end of an input or when the next character isn't a letter.
func (l *Lexer) readIdentifier(first byte) string {
    identifier := make([]byte, 0)
    identifier = append(identifier, first)

    for {
        peekedToken, _ := peekBuffer(*l.inputBuffer, 1)
        if isLetter(peekedToken) {
            next := l.inputBuffer.Next(1)
            identifier = append(identifier, next[0])
            l.currPosition++
            l.nextPosition++
        } else {
            break
        }
    }

    return string(identifier)
}

// peekBuffer peeks at the next byte.
// It's error should be ignored, since the only error should occur is io.EOF, which is also dealt with in ReadNextToken.
func peekBuffer(buf bytes.Buffer, n int) (byte, error) {
    if n <= 0 {
        return 0, nil
    }

    // Check if there are enough bytes in the buffer to peek.
    if buf.Len() < n {
        return 0, fmt.Errorf("buffer has less than %d bytes", n)
    }

    // Read the next n bytes without consuming them.
    peekedBytes, err := buf.ReadByte()

    if err != nil {
        // Including io.EOF and other unknown errors.
        return 0, err
    }

    return peekedBytes, nil
}

// isDigit determines whether an input character is a number.
func isDigit(ch byte) bool {
    // For ASCII, letter 0-9 lies between [48, 57].
    if 48 <= ch && ch <= 57 {
        return true
    }
    return false
}

// isLetter determines whether an input character is a letter.
func isLetter(ch byte) bool {
    // For ASCII, letter a-z lies within [97, 122] and A-Z within [65, 90].
    // We treat "_" as letter(ASCII: 95), indicating we allow both Camel Case and Snake Case for names of variables or functions.
    if (65 <= ch && ch <= 90) || (97 <= ch && ch <= 122) || ch == 95 {
        return true
    }
    return false
}

// isLetter determines whether an input character is a decimal point.
func isDecimalSeparator(ch byte) bool {
    return ch == 46
}

// isWhiteSpace determines whether an input character is a white space.
func isWhiteSpace(ch byte) bool {
    whiteSpaces := map[string]struct{}{" ": {}, "\n": {}, "\t": {}, "\r": {}}
    _, ok := whiteSpaces[string(ch)]
    return ok
}
