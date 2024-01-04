package token

import "testing"

func TestNew(t *testing.T) {
    testCase := struct {
        LexicalType string
        Literal     string
    }{INT, "5"}

    tok := New(testCase.LexicalType, testCase.Literal)
    if tok == nil {
        t.Errorf("Error created nil token.\n")
    }

    if tok.LexicalType != testCase.LexicalType {
        t.Errorf("Error token type: expected %s, got %s.\n", testCase.LexicalType, tok.LexicalType)
    }

    if tok.Literal != testCase.Literal {
        t.Errorf("Error token literal: expected %s, got %s.\n", testCase.Literal, tok.Literal)
    }
}
