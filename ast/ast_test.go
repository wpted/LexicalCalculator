package ast

import "testing"

func TestEvaluate(t *testing.T) {
    testCases := []struct {
        root   *Node
        left   *Node
        right  *Node
        result float32
    }{
        {
            root:   New(5, true, "", false, nil, nil),
            left:   nil,
            right:  nil,
            result: 5,
        },
    }

    for _, tc := range testCases {
        re, err := Evaluate(tc.root)
        if err != nil {
            t.Errorf("Error evaluating node, got err: %v.\n", err)
        }

        if re != tc.result {
            t.Errorf("Error evaluating node: expected %f, got %f.\n", tc.result, re)
        }
    }
}
