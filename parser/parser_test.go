package parser

import (
    "LexicalCalculator/ast"
    "LexicalCalculator/lexer"
    "LexicalCalculator/support"
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
            err   error
        }{
            // Error calculator prompt.
            {input: "", err: ErrPrompt},                  // missing the rest of the prompt.
            {input: "caldc '5 + 5'", err: ErrPrompt},     // misspelled calc.
            {input: "calc", err: ErrPrompt},              // missing the rest of the prompt.
            {input: "calc 5 + 5'", err: ErrOpeningQuote}, // missing first single quote.
            {input: "calc '5 + 5", err: ErrClosingQuote}, // missing closing single quote.

            // Error Opening brackets.
            {input: "calc '(0'", err: ErrEquation},
            {input: "calc '((0 + 1)'", err: ErrEquation},
            {input: "calc '[(12 + 3 * 6] + 1'", err: ErrEquation},
            {input: "calc '[(12 + 3) * 6 + 1'", err: ErrEquation},
            {input: "calc '{[(12 + 3) * 6] + 1 * 2'", err: ErrEquation},

            // Error closing brackets.
            {input: "calc '0)'", err: ErrEquation},
            {input: "calc '(0 + 1))'", err: ErrEquation},
            {input: "calc '[12 + 3 * 6)] + 1'", err: ErrEquation},
            {input: "calc '(12 + 3) * 6] + 1'", err: ErrEquation},
            {input: "calc '()1 * 2'", err: ErrEquation},
            {input: "calc '[(1]) * 2'", err: ErrEquation},

            // Error equation.
            {input: "calc ''", err: ErrEquation},
            {input: "calc '5 5'", err: ErrEquation},
            {input: "calc '5 +'", err: ErrEquation},

            {input: "calc '*5'", err: ErrEquation},
            {input: "calc '/5'", err: ErrEquation},
            {input: "calc '-5 +'", err: ErrEquation},
            // Zero division.
            {input: "calc '1 / 0'", err: ast.ErrZeroDivision},
        }

        for testNum, tc := range testCases {
            p.input(tc.input)
            err := p.parsePrompt()
            if testNum < 16 && !errors.Is(err, tc.err) {
                t.Errorf("error parsing incorrect prompt: expected error %s, got error %s.\n", tc.err, err)
            }

            // All equation error test cases are after testNum 55.
            if 17 <= testNum && testNum < 22 {
                _, err = p.parseEquation(0)
                if !errors.Is(err, tc.err) {
                    t.Errorf("error parsing incorrect equation: expected error %s, got error %s.\n", tc.err, err)
                }
            } else if testNum >= 22 {
                // For error equations like 1/0.
                n, _ := p.parseEquation(0)
                _, err = ast.Evaluate(n)
                if !errors.Is(err, tc.err) {
                    t.Errorf("error evaluating incorrect equation: expected error %s, got error %s.\n", tc.err, err)
                }
            }
        }
    })

    t.Run("Correct input", func(t *testing.T) {
        testCases := []struct {
            input  string
            tokens int
            result float32
            err    error
        }{
            // Binary operation.
            {input: "calc 'ans'", tokens: 1, result: 0},
            {input: "calc '0'", tokens: 1, result: 0},
            {input: "calc '1 + 2'", tokens: 3, result: 3},
            {input: "calc '2.1 - 1'", tokens: 3, result: 1.1},
            {input: "calc '2.2 * 3'", tokens: 3, result: 6.6},
            {input: "calc '3 / 2'", tokens: 3, result: 1.5},
            {input: "calc '2 + 3 + 4'", tokens: 5, result: 9},

            {input: "calc '2 + 3 * 5'", tokens: 5, result: 17},
            {input: "calc '1 + 2 * 3'", tokens: 5, result: 7},
            {input: "calc '3 * 7 + 5 * 4'", tokens: 7, result: 41},

            {input: "calc '2 + 2 + 2 + 2 - 1'", tokens: 9, result: 7},
            {input: "calc '2 + 2 + 2 + 2 - 1'", tokens: 9, result: 7},

            // Unary operation.
            {input: "calc '-ans'", tokens: 2, result: -7}, // The answer here should inherit the previous test result.
            {input: "calc '-1'", tokens: 2, result: -1},
            {input: "calc '-1 + 2'", tokens: 4, result: 1},
            {input: "calc '-1 + 2 * 5'", tokens: 6, result: 9},
            {input: "calc '--1'", tokens: 3, result: 1},
            {input: "calc '++5'", tokens: 3, result: 5},
            {input: "calc '5 + - 5'", tokens: 4, result: 0},

            // Brackets.
            {input: "calc '(0)'", tokens: 3, result: 0},
            {input: "calc '((0 + 1))'", tokens: 7, result: 1},
            {input: "calc '(1 + 2) * 3'", tokens: 7, result: 9},
            {input: "calc '(((0) + 1))'", tokens: 9, result: 1},
            {input: "calc '((2.1 + 3.7) * 4.425) * 6'", tokens: 11, result: 153.99},
            {input: "calc '(((1 + 2) * 3) + 4) * 5'", tokens: 15, result: 65},
            {input: "calc '((12 + 3) * (1 + 5)) + 1'", tokens: 15, result: 91},

            {input: "calc '[(12 + 3) * 6] + 1'", tokens: 11, result: 91},
            {input: "calc '[(1 + 2) * 3] * 4'", tokens: 11, result: 36},
            {input: "calc '[(ans + 8) * 2] / 8'", tokens: 11, result: 11},

            {input: "calc '{[(12 + 3) * 6] + 1} * 2'", tokens: 15, result: 182},
            {input: "calc '{[(1 + 2) * 3] + 4} * 5'", tokens: 15, result: 65},
            {input: "calc '{[(1 + 2) * 3] * [(100 / 20) + 8]} - 123'", tokens: 23, result: -6},
        }

        for _, tc := range testCases {
            p.input(tc.input)
            err := p.parsePrompt()
            if err != nil {
                if tc.err == nil {
                    t.Errorf("Error parsing calculator prompt, got error: %v.\n", err)
                } else {
                    if !errors.Is(err, tc.err) {
                        t.Errorf("error incorrect error: expected error %s, got error %s.\n", tc.err, err)
                    }
                }
            } else {
                // Check the first token, should be a calc token.
                if p.root.Token.LexicalType != token.CALC {
                    t.Errorf("error root token type: expected %s, got %s.\n", token.CALC, p.root.Token.LexicalType)
                }

                if p.root.Token.Literal != "calc" {
                    t.Errorf("error root token literal: expected %s, got %s.\n", "calc", p.root.Token.Literal)
                }

                if p.root.EquationTokens == nil {
                    t.Errorf("error parsing equation: expected non nil root.\n")
                }

                if len(p.root.EquationTokens) != tc.tokens {
                    t.Errorf("error parsing equation: incorrect amount of equation tokens, expected %d, got %d.\n", tc.tokens, len(p.root.EquationTokens))
                }

                n, err := p.parseEquation(0)
                if err != nil {
                    t.Errorf("Error parsing equation: got error %s.\n", err)
                }

                val, err := ast.Evaluate(n)
                if err != nil {
                    t.Errorf("error calculating: got error %s.\n", err)
                }
                p.result = val
                // Setting epsilon as accuracy.
                if !support.AlmostEqual(val, float64(tc.result), 0.0001) {
                    t.Errorf("error calculated value: expected %f, got %f.\n", tc.result, val)
                }
            }
        }
    })
}
