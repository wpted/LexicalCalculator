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
            _, err := p.Parse()
            if !errors.Is(err, i.error) {
                t.Errorf("error parsing incorrect prompt: expected error %v, got error %v.\n", i.error, err)
            }
        }
    })

    t.Run("Correct input", func(t *testing.T) {
        testCases := []struct {
            input  string
            result float32
        }{
            {input: "calc '1 + 2'", result: 3},
            {input: "calc '2 - 1'", result: 1},
            {input: "calc '2 * 3'", result: 6},
            {input: "calc '3 / 2'", result: 1.5},
        }

        for _, tc := range testCases {
            p.Input(tc.input)
            root, err := p.Parse()
            if err != nil {
                t.Errorf("Error parsing calculator prompt, got error: %v.\n", err)
            }

            // Check the first token, should be a calc token.
            if root.Token.LexicalType != token.CALC {
                t.Errorf("error root token type: expected %s, got %s.\n", token.CALC, root.Token.LexicalType)
            }

            if root.Token.Literal != "calc" {
                t.Errorf("error root token literal: expected %s, got %s.\n", "calc", root.Token.Literal)
            }

            // TODO: Check the parsed equation expression.
            if root.Equation == nil {
                t.Errorf("error parsing equation: expected non nil root")
            }

            result, _ := ast.Evaluate(root.Equation)
            if result != tc.result {
                t.Errorf("error evaluated result: expected %f, got %f.\n", tc.result, result)
            }
        }
    })
}
