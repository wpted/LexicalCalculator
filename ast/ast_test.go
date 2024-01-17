package ast

import (
    "testing"
)

func TestNode_String(t *testing.T) {
    testCases := []struct {
        node *Node
        str  string
    }{
        {
            node: New(nil, 1, true, "", false, nil, nil),
            str:  "1.0000",
        },
        {
            node: New(nil, 0, false, "+", true,
                nil,
                New(nil, 2, true, "", false, nil, nil),
            ),
            str: "(+ 0 2.0000)",
        },
        {
            node: New(nil, 0, false, "+", true,
                New(nil, 3, true, "", false, nil, nil),
                New(nil, 2, true, "", false, nil, nil),
            ),
            str: "(+ 3.0000 2.0000)",
        },
        {
            node: New(nil, 0, false, "-", true,
                New(nil, 3, true, "", false, nil, nil),
                New(nil, 2, true, "", false, nil, nil),
            ),
            str: "(- 3.0000 2.0000)",
        },
        {
            node: New(nil, 0, false, "*", true,
                New(nil, 3, true, "", false, nil, nil),
                New(nil, 2, true, "", false, nil, nil),
            ),
            str: "(* 3.0000 2.0000)",
        },
        {
            node: New(nil, 0, false, "/", true,
                New(nil, 3, true, "", false, nil, nil),
                New(nil, 2, true, "", false, nil, nil),
            ),
            str: "(/ 3.0000 2.0000)",
        },
        {
            node: New(nil, 0, false, "+", true,
                New(nil, 5, true, "", false, nil, nil),
                New(nil, 2, false, "*", true,
                    New(nil, 2, true, "", false, nil, nil),
                    New(nil, 3, true, "", false, nil, nil),
                ),
            ),
            str: "(+ 5.0000 (* 2.0000 3.0000))",
        },
        {
            node: New(nil, 0, false, "-", true,
                New(nil, 0, false, "+", true,
                    New(nil, 1, true, "", false, nil, nil),
                    New(nil, 0, false, "*", true,
                        New(nil, 2, true, "", false, nil, nil),
                        New(nil, 5, true, "", false, nil, nil),
                    ),
                ),
                New(nil, 0, false, "*", true,
                    New(nil, 3, true, "", false, nil, nil),
                    New(nil, 2, true, "", false, nil, nil),
                ),
            ),
            str: "(- (+ 1.0000 (* 2.0000 5.0000)) (* 3.0000 2.0000))",
        },
    }

    for _, tc := range testCases {
        sExpression := tc.node.String()
        if sExpression != tc.str {
            t.Errorf("Error transforming code into S-expression: expected %s, got %s.\n", tc.str, sExpression)
        }
    }
}

func TestEvaluate(t *testing.T) {
    leaf1 := New(nil, 6, true, "", false, nil, nil)
    leaf2 := New(nil, 3, true, "", false, nil, nil)
    leaf3 := New(nil, 5, true, "", false, nil, nil)
    leaf4 := New(nil, 0, false, "+", true, leaf1, leaf2)
    leaf5 := New(nil, 0, false, "*", true, leaf1, leaf2)
    testCases := []struct {
        root   *Node
        left   *Node
        right  *Node
        result float64
    }{
        {
            // nil node
            root:   nil,
            left:   nil,
            right:  nil,
            result: 0,
        },
        {
            // 6
            root:   leaf1,
            left:   nil,
            right:  nil,
            result: 6,
        },
        {
            // -6
            root:   New(nil, 0, false, "-", true, nil, leaf1),
            left:   nil,
            right:  leaf1,
            result: -6,
        },
        {
            // 6 + 3
            root:   New(nil, 0, false, "+", true, leaf1, leaf2),
            left:   leaf1,
            right:  leaf2,
            result: 9,
        },
        {
            // 6 - 3
            root:   New(nil, 0, false, "-", true, leaf1, leaf2),
            left:   leaf1,
            right:  leaf2,
            result: 3,
        },
        {
            // 6 * 3
            root:   New(nil, 0, false, "*", true, leaf1, leaf2),
            left:   leaf1,
            right:  leaf2,
            result: 18,
        },
        {
            // 6 / 3
            root:   New(nil, 0, false, "/", true, leaf1, leaf2),
            left:   leaf1,
            right:  leaf2,
            result: 2,
        },
        {
            // 6 / 3
            root:   New(nil, 0, false, "/", true, leaf1, leaf2),
            left:   leaf1,
            right:  leaf2,
            result: 2,
        },
        {
            // 5 + 6 + 3
            root:   New(nil, 0, false, "+", true, leaf3, leaf4),
            left:   leaf3,
            right:  leaf4,
            result: 14,
        },
        {
            // 5 + 6 * 3
            root:   New(nil, 0, false, "+", true, leaf3, leaf5),
            left:   leaf3,
            right:  leaf5,
            result: 23,
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
