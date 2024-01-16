package main

import (
    "LexicalCalculator/lexer"
    "LexicalCalculator/parser"
    "bufio"
    "fmt"
    "os"
    "strings"
)

const (
    REPL  = ">> "
    QUIT  = "quit"
    HELP  = "help"
    CLEAR = "clear"
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
    fmt.Println("Calculator started.")
    fmt.Println(">>>> Input your calculator prompt in format: calc '<your equation here>'")

    fmt.Println(">>>> Or type help to see instructions.")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Printf("%sInsert your prompt: ", REPL)
        scanner.Scan()
        cmd := scanner.Text()

        switch strings.ToLower(cmd) {
        case HELP:
            fmt.Println("Input prompts")
            fmt.Println("    - calc '<equation>'")
            fmt.Println("    - quit")
            fmt.Println("    - help")
        case CLEAR:
            p.Result = 0
            continue
        case QUIT:
            fmt.Println("Exit Calculator.")
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
