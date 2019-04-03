![Logo](logo.png)

### Intro
g-fu is a Lisp developed and embedded in Go. The initial [release](https://github.com/codr7/g-fu/tree/master/v1) weighs in at 1 kloc and implements a basic, extensible, tree-walking interpreter.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.3

Press Return twice to evaluate.

  (let (fib (fun (n)
              (if (< n 2)
                n
                (+ (fib (- n 1)) (fib (- n 2))))))
    (dump (fib 20)))

6765
```

### License
LGPL3