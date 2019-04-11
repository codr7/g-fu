![Logo](logo.png)

### Intro
g-fu is a pragmatic [Lisp](https://xkcd.com/297/) developed and embedded in Go. The initial [release](https://github.com/codr7/g-fu/tree/master/v1) implements an extensible, tree-walking interpreter with support for quotation and macros, symbols, lambdas, bindings, tail-call optimization and eval; all weighing in well below 2 kloc.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.5

Press Return twice to evaluate.

  (let (fib (fun (n)
              (if (< n 2)
                n
                (+ (fib (- n 1)) (fib (- n 2))))))
    (dump (fib 20)))

6765
```

### Macros
One of the most common macro examples is the `while`-loop. The example below defines it in terms of a more general `loop`-macro, which will follow shortly. Note that g-fu uses `%` as opposed to `,` for interpolating values, `_` in place of `nil` and `..` to splat.

```
  (let while (macro (cond body..)
    '(loop
       (if %cond _ (break))
       %body..)))

  (let (i 0)
    (while (< i 7)
      (dump i)
      (inc i)))

0
1
2
3
4
5
6
```

`loop` supports exiting with a result using `break` within its body, which is trapped by a nested macro. Most of the hard work is performed by an anonymous, tail-recursive function; fresh argument symbols are created to avoid capturing the calling environment.

```
  (let loop (macro (body..)
    (let done (g-sym) result (g-sym))
  
    '(let (break (macro (args..) '(recall T %args..)))
       ((fun (%done? %result..)
          (if %done %result.. (do %body.. (recall))))))))

  (dump (loop (dump 'foo) (break 'bar) (dump 'baz)))

'foo
'bar
```

### Profiling
CPU profiling may be enabled by passing `-prof` on the command line; results are written to the specified file, `fib_tail.prof` in the following example.

```
$ ./gfu -prof fib_tail.prof bench/fib_tail.gf

$ go tool pprof fib_tail.prof
File: gfu
Type: cpu
Time: Apr 12, 2019 at 12:52am (CEST)
Duration: 16.52s, Total samples = 17.31s (104.79%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 10890ms, 62.91% of 17310ms total
Dropped 99 nodes (cum <= 86.55ms)
Showing top 10 nodes out of 79
      flat  flat%   sum%        cum   cum%
    2130ms 12.31% 12.31%    15980ms 92.32%  _/home/a/Dev/g-fu/v1/src/gfu.Vec.EvalVec
    1580ms  9.13% 21.43%    15980ms 92.32%  _/home/a/Dev/g-fu/v1/src/gfu.(*VecType).Eval
    1500ms  8.67% 30.10%     1680ms  9.71%  runtime.heapBitsSetType
    1180ms  6.82% 36.92%     1520ms  8.78%  _/home/a/Dev/g-fu/v1/src/gfu.(*Env).Find
     970ms  5.60% 42.52%     2290ms 13.23%  _/home/a/Dev/g-fu/v1/src/gfu.(*SymType).Eval
     830ms  4.79% 47.31%    15980ms 92.32%  _/home/a/Dev/g-fu/v1/src/gfu.Val.Eval
     780ms  4.51% 51.82%     4440ms 25.65%  runtime.mallocgc
     770ms  4.45% 56.27%      770ms  4.45%  runtime.memclrNoHeapPointers
     730ms  4.22% 60.49%    15980ms 92.32%  _/home/a/Dev/g-fu/v1/src/gfu.(*FunType).Call
     420ms  2.43% 62.91%      420ms  2.43%  runtime.memmove
```

### License
LGPL3