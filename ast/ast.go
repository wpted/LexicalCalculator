package ast

import (
    "LexicalCalculator/token"
    "errors"
)

// an equation like this
// 5 + 2 * 3
// should be represented as
//       +
//      / \
//     5   *
//        / \
//       2   3

var (
    ErrZeroDivision = errors.New("error cannot use 0 as denominator")
)

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
    Value      float32
    Left       *Node
    Right      *Node
}

// New creates a new Node.
func New(value float32, isValue bool, operation string, isOperator bool, leftChild *Node, rightNode *Node) *Node {
    return &Node{
        IsOperator: isOperator,
        Operator:   operation,
        IsValue:    isValue,
        Value:      value,
        Left:       leftChild,
        Right:      rightNode,
    }
}

// Evaluate evaluates the current node and return the result of the equation.
func Evaluate(equationNode *Node) (float32, error) {
    // Scenarios
    // 1. 6 ( One single integer )
    // 2. -6 ( Negative integer )
    // 3. 6 + 5 ( Binary Expression )
    // 4. nil node ( should be 0 )
    // 5. 5 + 2 * 3 ( '*' operator has higher precedence than '+' )

    // If node is neither an operator nor a value, it is 0.
    // Consider ' + 5 '
    // it's   +
    //       / \
    //     nil  5
    // the value of the left child should be considered as 0.
    if equationNode == nil {
        return 0, nil
    }

    if equationNode.IsValue {
        return equationNode.Value, nil
    }

    if equationNode.IsOperator {
        left, err := Evaluate(equationNode.Left)
        right, err := Evaluate(equationNode.Right)

        switch equationNode.Operator {
        case "+":
            return left + right, err
        case "-":
            return left - right, err
        case "*":
            return left * right, err
        case "/":
            if right == 0 {
                return 0, ErrZeroDivision
            }
            return left / right, err
        }
    }

    return 0, nil
}
