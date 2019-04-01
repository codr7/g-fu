![Logo](logo.png)

### Intro
g-fu is a Lisp developed and embedded in Go. The implementation is divided into stages, each adressing a specific issue and culminating in a separate dialect that may be considered frozen in time except for the occasional bug fix / refactoring. The initial [dialect](https://github.com/codr7/g-fu/tree/master/v1) weighs in at 1 kloc and implements a tree-walking interpreter to the point where it can run a REPL and calculate the Fibonacci sequence using tail-recursion.

### License
AGPL