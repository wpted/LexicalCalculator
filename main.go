package main

import (
    "LexicalCalculator/lexer"
    "LexicalCalculator/parser"
    "bufio"
    "fmt"
    "os"
)

const (
    REPL = ">> "
    QUIT = "quit"
)

func main() {
    // Initialize lexer and parser.
    l := lexer.New()
    p := parser.New(l)

    // We want a calculator that reads scripts like `calc "1 + 1"`, `calc "2 * 3 + 4"` or `quit`.
    // Operators: +, -, *, /, **, (), [], {}.
    // Should have prefix calculation like `calc "-5 + 4"`
    // Should handle error like zero-divisions.
    // Create CLI environment(repl), takes commands; ['calc', 'end'].
    fmt.Println("Start calculator.")
    fmt.Println(">>>>")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print(REPL)
        scanner.Scan()
        cmd := scanner.Text()

        switch cmd {
        case QUIT:
            return
        default:
            calculatedResult, err := p.Evaluate(cmd)
            if err != nil {
                fmt.Printf("%sIncorrect prompt: %s\n", REPL, cmd)
                continue
            }
            // Round the result to 2 decimal places
            fmt.Printf(">> result: %.2f\n", calculatedResult)
        }
    }
}
