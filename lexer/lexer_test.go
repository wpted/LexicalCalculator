package lexer

import (
    "LexicalCalculator/token"
    "testing"
)

func TestLexer_Input(t *testing.T) {
    l := New()
    testCases := []string{"Hello, world!", "Goodbye, world!"}

    for _, tc := range testCases {
        l.Input(tc)

        if l.inputBuffer.String() != tc {
            t.Errorf("Error writing data into lexer: expected %s, got %s.\n", tc, l.inputBuffer.String())
        }
    }
}

func TestLexer_ReadNextToken(t *testing.T) {
    l := New()
    testCases := []struct {
        input  string
        result []token.Token
    }{
        {
            input: "",
            result: []token.Token{
                {Literal: token.EOF, LexicalType: token.EOF},
            },
        },
        {
            input: "+-*/'()[]{}",
            result: []token.Token{
                {Literal: "+", LexicalType: token.PLUS},
                {Literal: "-", LexicalType: token.MINUS},
                {Literal: "*", LexicalType: token.ASTERISK},
                {Literal: "/", LexicalType: token.SLASH},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "(", LexicalType: token.LPAREN},
                {Literal: ")", LexicalType: token.RPAREN},
                {Literal: "[", LexicalType: token.LSQBRACK},
                {Literal: "]", LexicalType: token.RSQBRACK},
                {Literal: "{", LexicalType: token.LCURBRACK},
                {Literal: "}", LexicalType: token.RCURBRACK},
                {Literal: token.EOF, LexicalType: token.EOF},
            },
        },
        {
            input: "    - \n+\r-*/\t'",
            result: []token.Token{
                {Literal: "-", LexicalType: token.MINUS},
                {Literal: "+", LexicalType: token.PLUS},
                {Literal: "-", LexicalType: token.MINUS},
                {Literal: "*", LexicalType: token.ASTERISK},
                {Literal: "/", LexicalType: token.SLASH},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: token.EOF, LexicalType: token.EOF},
            },
        },
        {
            input: "368000 12345",
            result: []token.Token{
                {Literal: "368000", LexicalType: token.INT},
                {Literal: "12345", LexicalType: token.INT},
                {Literal: token.EOF, LexicalType: token.EOF},
            },
        },
        {
            input: "3..68 1.",
            result: []token.Token{
                {Literal: "3..68", LexicalType: token.UNKNOWN},
                {Literal: "1.", LexicalType: token.UNKNOWN},
                {Literal: token.EOF, LexicalType: token.EOF},
            },
        },
        {
            input: "calc ans hello",
            result: []token.Token{
                {Literal: "calc", LexicalType: token.CALC},
                {Literal: "ans", LexicalType: token.ANS},
                {Literal: "hello", LexicalType: token.UNKNOWN},
                {Literal: token.EOF, LexicalType: token.EOF},
            },
        },
        {
            input: "calc '5 + 56'",
            result: []token.Token{
                {Literal: "calc", LexicalType: token.CALC},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "5", LexicalType: token.INT},
                {Literal: "+", LexicalType: token.PLUS},
                {Literal: "56", LexicalType: token.INT},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "EOF", LexicalType: token.EOF},
            },
        },
        {
            input: "calc '2 * 4'",
            result: []token.Token{
                {Literal: "calc", LexicalType: token.CALC},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "2", LexicalType: token.INT},
                {Literal: "*", LexicalType: token.ASTERISK},
                {Literal: "4", LexicalType: token.INT},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "EOF", LexicalType: token.EOF},
            },
        },
        {
            input: "calc '4 / 2'",
            result: []token.Token{
                {Literal: "calc", LexicalType: token.CALC},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "4", LexicalType: token.INT},
                {Literal: "/", LexicalType: token.SLASH},
                {Literal: "2", LexicalType: token.INT},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "EOF", LexicalType: token.EOF},
            },
        },
        {
            input: "calc '4 - 2'",
            result: []token.Token{
                {Literal: "calc", LexicalType: token.CALC},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "4", LexicalType: token.INT},
                {Literal: "-", LexicalType: token.MINUS},
                {Literal: "2", LexicalType: token.INT},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "EOF", LexicalType: token.EOF},
            },
        },
        {
            input: "calc '2.1 + 3.5'",
            result: []token.Token{
                {Literal: "calc", LexicalType: token.CALC},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "2.1", LexicalType: token.FLOAT},
                {Literal: "+", LexicalType: token.PLUS},
                {Literal: "3.5", LexicalType: token.FLOAT},
                {Literal: "'", LexicalType: token.SINGLEQUOTE},
                {Literal: "EOF", LexicalType: token.EOF},
            },
        },
    }

    for _, tc := range testCases {
        l.Input(tc.input)
        var currTokenPosition int
        for {
            tok := l.ReadNextToken()
            expectedToken := tc.result[currTokenPosition]
            if tok.Literal != expectedToken.Literal {
                t.Errorf("Error token literal: expected %s, got %s.\n", expectedToken.Literal, tok.Literal)
            }

            if tok.LexicalType != expectedToken.LexicalType {
                t.Errorf("Error token type: expected %s, got %s.\n", expectedToken.LexicalType, tok.LexicalType)
            }

            if tok.LexicalType == token.EOF {
                break
            }
            currTokenPosition++
        }
        if currTokenPosition+1 != len(tc.result) {
            t.Errorf("Error parsed token length, expected %d, got %d.\n", len(tc.result), currTokenPosition+1)
        }
    }
}
