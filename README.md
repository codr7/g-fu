![Logo](logo.png)

### Intro
g-fu is a Lisp developed and embedded in Go. The first [release](https://github.com/codr7/g-fu/tree/master/v1) weighs in at 1 kloc and implements a simple, extensible tree-walking interpreter to the point where it's capable of calculating the Fibonacci sequence using tail-recursion.

```
  (let (fib (fun (n a b)
              (if n 
                (if (= n 1)
                  b
                  (recall (- n 1) b (+ a b)))
                a)))
    (dump (fib 20 0 1)))

6765
```

### License
AGPL