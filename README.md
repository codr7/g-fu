![Logo](logo.png)
  
### Intro
g-fu is a pragmatic [Lisp](https://xkcd.com/297/) developed and embedded in Go.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.12

Press Return twice to evaluate.

  (fun fib (n)
    (if (< n 2)
      n
      (+ (fib (- n 1)) (fib (- n 2)))))
```
```
  (fib 20)

6765
```

### Goals
The primary goal is to provide a fully integrated, practical Lisp in Go. Practical as a complement to Go; which means a clean implementation that composes well with Go is more important than raw performance, among other things.

### Status
The initial [release](https://github.com/codr7/g-fu/tree/master/v1) is more or less ready for action. Once it stabilizes; work begins on v2, which will focus on compilation to an internal representation optimized for evaluation.

### License
LGPL3