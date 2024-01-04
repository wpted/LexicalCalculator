/*
Package lexer implements utilities of a lexical parser.
Lexer takes input and parse input into token ( if the lexer understands the input ).
Should be able to reuse, not creating a lexer everytime we receive an input.
*/
package lexer

import (
    "LexicalCalculator/token"
    "bytes"
)

// Lexer is the lexical analyzer used in the calculator.
// It's supposed to be called by a parser and will lazily return the next token on the fly.
// Lexer resets inputBuffer and resets field currPosition and nextPosition after receiving a new input.
type Lexer struct {
    inputBuffer  bytes.Buffer
    currPosition int
    nextPosition int
    errors       []error
}

// New creates a new Lexer.
func New() *Lexer {
    return &Lexer{
        inputBuffer:  bytes.Buffer{},
        currPosition: -1,
        nextPosition: 0,
    }
}

// Input writes data into the lexer. If the buffer is not empty, it frees the buffer first.
func (l *Lexer) Input(data string) {
    l.free()
    // Write appends the data to the buffer, growing the buffer as needed. The return value n is the length of p; err is always nil.
    // If the buffer becomes too large, Write will panic with ErrTooLarge.
    _, _ = l.inputBuffer.Write([]byte(data))
}

// free resets the buffer to be empty, but it retains the underlying storage for use by future writes.
func (l *Lexer) free() {
    l.inputBuffer.Reset()
    l.currPosition = -1
    l.nextPosition = 0
}

// ReadNextToken returns a token after every read, skipping all encountered white spaces.
func (l *Lexer) ReadNextToken() *token.Token {
    var tok *token.Token
    if l.nextPosition > l.inputBuffer.Len() {
        tok = token.New("EOF", token.EOF)
    }

    next := l.inputBuffer.Next(1)
    l.currPosition++
    l.nextPosition++

    // When encounter a white space, advance the pointer.
    whiteSpaces := map[string]struct{}{
        " ":  {},
        "\n": {},
        "\t": {},
        "\r": {},
    }

    _, ok := whiteSpaces[string(next)]
    for ok {
        next = l.inputBuffer.Next(1)
        l.currPosition++
        l.nextPosition++
        _, ok = whiteSpaces[string(next)]
    }

    switch string(next) {
    case "":
        return tok
    case "'":
        tok = token.New(token.SINGLEQUOTE, string(next))
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
            // What we do here is to advance the pointer and fetch entire token.
            literal := l.readNumber(next[0])
            tok = token.New(token.INT, literal)
        } else if isLetter(next[0]) {
            literal := l.readIdentifier(next[0])
            if literal == "calc" {
                tok = token.New(token.CALC, literal)
            } else {
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
func (l *Lexer) readNumber(first byte) string {
    digits := make([]byte, 0)
    digits = append(digits, first)
    next := l.inputBuffer.Next(1)

    // Check if there's more after the first digit.
    if len(next) != 0 {
        for isDigit(next[0]) {

            digits = append(digits, next[0])

            l.currPosition++
            l.nextPosition++
            next = l.inputBuffer.Next(1)
            if len(next) == 0 {
                // If the next read byte is empty, we've reached an end.
                // break here to prevent index out of range error (since next is []byte{}, index 0 is out of range).
                break
            }
        }
    }

    return string(digits)
}

// readNumber returns the literal of an identifier.
// It advances the pointer until it's at the end of an input or when the next character isn't a letter.
func (l *Lexer) readIdentifier(first byte) string {
    identifier := make([]byte, 0)
    identifier = append(identifier, first)

    next := l.inputBuffer.Next(1)

    if len(next) != 0 {
        for isLetter(next[0]) {

            identifier = append(identifier, next[0])

            l.currPosition++
            l.nextPosition++
            next = l.inputBuffer.Next(1)
            if len(next) == 0 {
                break
            }
        }
    }

    return string(identifier)
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
