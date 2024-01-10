package parser

import (
    "LexicalCalculator/ast"
    "LexicalCalculator/lexer"
    "LexicalCalculator/token"
    "errors"
    "testing"
)

func TestParser_Parse(t *testing.T) {
    l := lexer.New()
    p := New(l)

    t.Run("Incorrect input", func(t *testing.T) {
        testCases := []struct {
            input string
            error error
        }{
            {"", ErrPrompt},                  // missing the rest of the prompt.
            {"caldc '5 + 5'", ErrPrompt},     // misspelled calc.
            {"calc", ErrPrompt},              // missing the rest of the prompt.
            {"calc 5 + 5'", ErrOpeningQuote}, // missing first single quote.
            {"calc '5 + 5", ErrClosingQuote}, // missing closing single quote.
        }

        for _, i := range testCases {
            p.Input(i.input)
            err := p.Parse()
            if !errors.Is(err, i.error) {
                t.Errorf("error parsing incorrect prompt: expected error %v, got error %v.\n", i.error, err)
            }
        }
    })

    t.Run("Correct input", func(t *testing.T) {
        testCases := []struct {
            input  string
            tokens int
            result float32
        }{
            {input: "calc '1'", tokens: 1, result: 1},
            {input: "calc '1 + 2'", tokens: 3, result: 3},
            {input: "calc '2 - 1'", tokens: 3, result: 1},
            {input: "calc '2 * 3'", tokens: 3, result: 6},
            {input: "calc '3 / 2'", tokens: 3, result: 1.5},
            {input: "calc '2 + 3 + 4'", tokens: 5, result: 9},
            {input: "calc '2 + 3 * 5'", tokens: 5, result: 17},
            {input: "calc '2 + 2 + 2 + 2 - 1'", tokens: 9, result: 7},
        }

        for _, tc := range testCases {
            p.Input(tc.input)
            err := p.Parse()
            if err != nil {
                t.Errorf("Error parsing calculator prompt, got error: %v.\n", err)
            }

            // Check the first token, should be a calc token.
            if p.root.Token.LexicalType != token.CALC {
                t.Errorf("error root token type: expected %s, got %s.\n", token.CALC, p.root.Token.LexicalType)
            }

            if p.root.Token.Literal != "calc" {
                t.Errorf("error root token literal: expected %s, got %s.\n", "calc", p.root.Token.Literal)
            }

            if p.root.EquationTokens == nil {
                t.Errorf("error parsing equation: expected non nil root")
            }

            if len(p.root.EquationTokens) != tc.tokens {
                t.Errorf("error parsing equation: incorrect amount of equation tokens, expected %d, got %d.\n", tc.tokens, len(p.root.EquationTokens))
            }

            n := p.parseEquation(0)
            val, _ := ast.Evaluate(n)
            if val != tc.result {
                t.Errorf("error calculated value: expected %f, got %f.\n", tc.result, val)
            }
        }
    })
}
