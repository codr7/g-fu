![Logo](logo.png)
  
### Intro
g-fu is a pragmatic [Lisp](https://xkcd.com/297/) developed and embedded in Go.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.15

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

### Status
The initial [release](https://github.com/codr7/g-fu/tree/master/v1) is more or less ready for action. Once it stabilizes; work begins on v2, which will focus on compilation to an internal representation optimized for evaluation.

### License
LGPL3