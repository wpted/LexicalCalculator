# Write Your Own Calculator

This is the solution for [Calculator Challenge](https://codingchallenges.fyi/challenges/challenge-calculator)
implemented using Go.

## Idea

Parsing is the process by which a compiler turns a sequence of tokens into a tree representation:

```go
                            Add
                 Parser     / \
'1 + 2 * 3'    ------->    1  Mul
                              / \
                             2   3
```

We are to take prompts like:

```go
    calc '1 + 2 * 3'
```

and return the result.

Operation precedence are taken cared using binding power like below:

```text
  
Equation            -   1     +     3   *   2   ^   2
                   / \       / \       / \     / \
Binding power         5     1   2     3   4   6   7
                     rbp   lbp  rbp
```

In each loop we compare lbp and the passed-in min_bp. If lbp is greater, the operator has higher priority.


## Implementation Steps

1. Create a lexer that tokenizes the input.
2. Create a parser that parses the tokens, and turn the tokens into an AST.
   To build the AST, I selected **Pratt Parsing** and **S-Expression** instead of using the Shunting Yard algorithm and stacks.
3. Evaluate the AST and output the calculated results.

## Program

Git pull the repo and run the program.

```shell
   go run main.go
```

or compile it first then run the corresponding executable on different platforms.

```shell
  # On Linux or MacOS
  go build -o calculator
  ./calculator
  
  # On Windows
  env GOOS=windows GOARCH=amd64 go build -o calculator.exe
  calculator.exe
```

You can also run the existing binaries pulled directly from the repo.

```shell
  # On Linux
  
  ./calculator
  
  # On Windows
  calculator.exe
```

After starting the program you'll see:

```go
    Calculator Started.
    >>>> Input your calculator prompt in format - calc '<your equation here>'
    >>>> Or type help to see instructions.
    >> Insert your prompt:
```

There are two types of prompt:

1. **Calculator prompts**
2. **Quit prompts**

### Prompt Example:

This is the general input format the calculator will take:

```go
   // 1. Calculator prompt, white space ignored, case-insensitive.
   calc '5 + 2 * 3'
```

To quit the calculator:

```go
  // 2. Quit prompt, case insensitive.    
  quit  
```

### Features
- [x] Float supports
    ```go
        calc '2.1 * 3.5'
    ```

- [x] Four simple expressions:

  ```go
    calc '1 + 2' // result: 3.0000
    calc '2 - 1' // result: 1.0000
    calc '2 * 3' // result: 6.0000
    calc '3 / 2' // result: 1.5000
  ```

- [x] Power with **^**
  
  ```go
    // Power with integers 
    calc '2 ^ 3'        // result: 6.0000
    calc '5.5 ^ 6'      // result: 27680.6406
    calc '-1 ^ 4'       // result: -1.0000
    calc '(-1) ^ 4'     // result: 1.0000
  ```

- [x] Mixed operations:

  ```go
    calc '1 + 2 * 3'             // result: 7.0000
    calc '3 * 7 + 5 * 4'         // result: 41.0000
    calc '3 * 7 + 2 ^ 4 - 7 ^ 0' // result: 36.0000
  ```

  The result is rounded to 4 decimal places.

- [x] Brackets:

  ```go
    calc '(1 + 2) * 3'                              // result: 9.0000
    calc '[(1 + 2) * 3] * 4'                        // result: 36.0000
    calc '{[(1 + 2) * 3] + 4} * 5'                  // result: 65.0000
    calc '{[(1 + 2) * 3] * [(100 / 20) + 8]} - 123' // result: -6.0000

    // Also supports equations with all parenthesis.
    calc '(((1 + 2) * 3) + 4) * 5'                  // result: 65.0000
  ```
  
  Mixed operations of brackets and **^** are handled.

  ```go
    // Should be different
    calc '-1 ^ 4'   // result: -1.0000
    calc '(-1) ^ 4' // result: 1.0000
  ```


  - Bracket expressions like below should cause error (unbalanced brackets):

    ```go
        calc '(1 + 2 * 3'
        calc '1 + 2) * 3'

        calc '[(1 + 2 * 3] * 4'
        calc '[1 + 2) * 3] * 4'
        calc '[(1 + 2) * 3 * 4'

        calc '{1 + 2) * 3] + 4} * 5'
    ```

- [x] Store previous result in **ans**.
  
  ```go
    // First prompt
    calc '1 + 2'        // result: 3.0000
        
    // The current 'ans' is 36.00
    calc 'ans'          // result: 36.0000
            
    // use 'ans' to called stored results
    calc 'ans * 12'     // result: 36.0000
    ```
- [x] Clear (AC button).
  
  ```go
    // First prompt
    calc '1 + 2'        // result: 3.0000
   
    // use 'ans' to call stored results
    calc 'ans * 12'     // result: 36.0000
   
    // The current 'ans' is 36.00
    calc 'ans'          // result: 36.0000
   
    // Clear 'ans'
    clear
   
    // The current 'ans' is 0.0000
    calc 'ans'          // result: 0.0000 
    ```

## References

### Tools:

- [AST Explorer](https://astexplorer.net)

### Video:

- [Swift 3 Fun Algorithms: Abstract Syntax Tree](https://www.youtube.com/watch?v=r14Vtwi2k7s)

### Reads:

- [Wikipedia: Shunting yard algorithm](https://en.wikipedia.org/wiki/Shunting_yard_algorithm#The_algorithm_in_detail)
- [Wikipedia: S-expression](https://en.wikipedia.org/wiki/S-expression)
- [Simple but Powerful Pratt Parsing](https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html#From-Precedence-to-Binding-Power)
- [How Desmos uses Pratt Parsers](https://engineering.desmos.com/articles/pratt-parser/)
- [Pratt Parsers: Expression Parsing Made Easy](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/)
- [Arrow functions break JavaScript parsers](https://dev.to/samthor/arrow-functions-break-javascript-parsers-1ldp)
- [手寫一個Parser - 代碼簡單而功能強大的Pratt Parsing](https://zhuanlan.zhihu.com/p/471075848)
