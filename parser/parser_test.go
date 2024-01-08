package parser

import (
    "LexicalCalculator/lexer"
    "LexicalCalculator/token"
    "errors"
    "testing"
)

func TestParser_Parse(t *testing.T) {
    l := lexer.New()
    p := New(l)

    t.Run("Correct input", func(t *testing.T) {
        input := "calc '5 + 5'"
        p.l.Input(input)
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
            t.Errorf("error parsing equation")
        }

    })

    t.Run("Incorrect input", func(t *testing.T) {
        testCases := []struct {
            input string
            error error
        }{
            {"caldc '5 + 5'", ErrPrompt},     // misspelled calc
            {"calc 5 + 5'", ErrOpeningQuote}, // missing first single quote
            {"calc '5 + 5", ErrClosingQuote}, // missing closing single quote
        }

        for _, i := range testCases {
            p.l.Input(i.input)
            _, err := p.Parse()
            if !errors.Is(err, i.error) {
                t.Errorf("error parsing incorrect prompt: expected error %v, got error %v.\n", i.error, err)
            }
        }
    })
}
