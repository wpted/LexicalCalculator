package ast

import "LexicalCalculator/token"

// an equation like this
// 5 + 2 * 3
// should be represented as
//       +
//      / \
//     5   *
//        / \
//       2   3

// Root should be the ast root of a calculator prompt.
type Root struct {
    Token    *token.Token
    Equation *Node
}

// Node is the general structure of all expressions in the equation.
type Node struct {
    Token      *token.Token
    IsOperator bool
    Operator   string
    IsValue    bool
    Value      int
    Left       *Node
    Right      *Node
}
