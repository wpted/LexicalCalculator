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

## Implementation Steps

1. Create a lexer that tokenizes the input.
2. Create a parser that parses the tokens, and turns it into an AST.
3. Evaluate the AST and output the calculated results.

## Current Feature

Run the program.

```shell
  go run main.go
```

After starting the program you'll see:

```go
Calculator Started.
>>>> Input your calculator prompt in format - calc '<your equation here>'
>>>> Or type help to see instructions.
>> Insert your prompt:
```

There is two type of prompt:

1. **Calculator prompts**
2. **Quit prompts**

### Prompt Example:

This is the general input format the calculator will take:

```go
   // 1. Calculator prompt, white space ignored, case insensitive.
   calc '5 + 2 * 3'
```

To quit the calculator:

```go
  // 2. Quit prompt, case insensitive.    
  quit  
```

### Supports

The calculator currently supports four simple expressions:

```go
    calc '1 + 2' // result: 3.00
    calc '2 - 1' // result: 1.00
    calc '2 * 3' // result: 6.00
    calc '3 / 2' // result: 1.50
```

and mixed operations like:

```go
    calc '1 + 2 * 3'     // result: 7.00
    calc '3 * 7 + 5 * 4' // result: 41.00
```

The result is rounded to 2 decimal places.

## TODOs

1. Float supports
    ```go
        calc '2.1 * 3.5'
    ```
   
2. Implement brackets: ( )[ ]{ }.
    ```go
        calc '{[(1 + 2) * 3] * [(100 / 20) + 8]} - 123'
    ```
   
3. Power with integers, sin, cos, tan.
    ```go
        // Power with integers
        calc '2 ^ 3'
        calc '5.5 ^ 6' 
    ```
    ```go
        // sin, cos, tan
        calc 'sin(37)' // result: 0.80
        calc 'cos(37)' // result: 0.80 
        calc 'tan(37)' // result: 0.75
    ```
   
4. Store previous result in **ans**.
    ```go
        // First prompt
        calc '1 + 2'        // result: 3.00
        
        // The current 'ans' is 36.00
        calc 'ans'          // result: 36.00
            
        // use 'ans' to called stored results
        calc 'ans * 12'     // result: 36.00 
    ```
5. Clear (AC button).
    ```go
        // First prompt
        calc '1 + 2'        // result: 3.00
   
        // use 'ans' to called stored results
        calc 'ans * 12'     // result: 36.00
   
        // The current ans is 36.00
        calc 'ans'          // result: 36.00
   
        // Clear 'ans'
        clear
   
        // The current ans is 0.00
        calc 'ans'          // result: 0.00 
    ```
    

## References

- Tools:
    - [AST Explorer](https://astexplorer.net)
- Video:
    - [Swift 3 Fun Algorithms: Abstract Syntax Tree](https://www.youtube.com/watch?v=r14Vtwi2k7s)
- Reads:
    - [S-expression according to Wikipedia](https://en.wikipedia.org/wiki/S-expression)
    - [Simple but Powerful Pratt Parsing](https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html#From-Precedence-to-Binding-Power)