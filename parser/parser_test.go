package parser

import (
    "LexicalCalculator/lexer"
    "errors"
    "testing"
)

func TestParser_Parse(t *testing.T) {
    l := lexer.New()
    p := New(l)

    t.Run("Correct input", func(t *testing.T) {
        input := "calc '5 + 5'"
        p.l.Input(input)
        _, err := p.Parse()
        if err != nil {
            t.Errorf("")
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
