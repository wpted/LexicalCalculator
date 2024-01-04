package main

import (
    "bufio"
    "fmt"
    "os"
)

const (
    QUIT = "quit"
)

func main() {
    // We want a calculator that reads scripts like `calc "1 + 1"`, `calc "2 * 3 + 4"` or `quit`.
    // Operators: +, -, *, /, **, (), [], {}.
    // Should have prefix calculation like `calc "-5 + 4"`
    // Should handle error like zero-divisions.
    // TODO 0: Create CLI environment(repl), takes commands; ['calc', 'end'].
    fmt.Println("Start calculator.")
    fmt.Println(">>>>")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print(">> ")
        scanner.Scan()
        cmd := scanner.Text()

        switch cmd {
        case QUIT:
            return
        default:

            fmt.Printf(">> result: %s\n", cmd)
        }
    }

    // TODO 1: Read expressions +, -, *, /.
    // Handle four simple expressions.
}
