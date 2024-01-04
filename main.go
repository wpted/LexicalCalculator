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
    // We want a calculator that reads scripts like `calc '1 + 1'`, `calc '2 * 3 + 4'` or `end`.
    // Operators: +, -, *, /, **, (), [], {}.
    // Should handle error like zero-divisions.
    // TODO 0: Create CLI environment(repl), takes commands; ['calc', 'end'].
    fmt.Println("Start calculator.")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        scanned := scanner.Scan()
        if !scanned {
            return
        }
        cmd := scanner.Text()
        if cmd == QUIT {
            break
        } else {
            fmt.Println(cmd)
        }

    }

    // TODO 1: Read expressions +, -, *, /.
}
