![Logo](logo.png)

### Intro
g-fu is a Lisp developed and embedded in Go. The initial [release](https://github.com/codr7/g-fu/tree/master/v1) weighs in at 1 kloc and implements a basic, extensible, tree-walking interpreter.

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
One of the most common macro examples is the `while`-loop. The example below defines it in terms of a more general `loop`-macro, which will follow shortly. Note that g-fu uses `%` as opposed to `,` for interpolating values, `_` in place of `nil` and `..` to splat values.

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

`loop` allows exiting with a result by calling `break` anywhere within the body. Most of the hard work is performed by an anonymous, tail-recursive function. A locally scoped macro is used to trap `break`-calls and a fresh symbol is allocated for the variable `break-args` to prevent potentially capturing the calling environment.

```
  (let loop (macro (body..)
         (let break-args (Sym))
         '(let (break (macro (args..)
                  '(recall %args..)))
            ((fun (%break-args)
               (or %break-args (do %body.. (recall _))))
             _))))

  (dump (loop (dump 'foo) (break 'bar) (dump 'baz)))

'foo
'bar
```

### Profiling
CPU profiling may be enabled by passing `-prof` with a filename on the command line.

```
$ ./gfu -prof fib_rec.prof bench/fib_rec.gf

$ go tool pprof fib_rec.prof 
File: gfu
Type: cpu
Time: Apr 6, 2019 at 4:29pm (CEST)
Duration: 16.52s, Total samples = 19.72s (119.38%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 12320ms, 62.47% of 19720ms total
Dropped 124 nodes (cum <= 98.60ms)
Showing top 10 nodes out of 95
      flat  flat%   sum%        cum   cum%
    3410ms 17.29% 17.29%     3650ms 18.51%  runtime.heapBitsSetType
    2370ms 12.02% 29.31%     2370ms 12.02%  runtime.memclrNoHeapPointers
    1520ms  7.71% 37.02%     1520ms  7.71%  runtime.memmove
    1100ms  5.58% 42.60%     2010ms 10.19%  runtime.scanobject
    1000ms  5.07% 47.67%    10460ms 53.04%  runtime.mallocgc
     920ms  4.67% 52.33%     1200ms  6.09%  _/home/a/Dev/g-fu/v1/src/gfu.(*Env).Find
     530ms  2.69% 55.02%    16980ms 86.11%  _/home/a/Dev/g-fu/v1/src/gfu.(*ExprForm).Eval
     530ms  2.69% 57.71%    16980ms 86.11%  _/home/a/Dev/g-fu/v1/src/gfu.VecForm.Eval
     490ms  2.48% 60.19%     9200ms 46.65%  runtime.growslice
     450ms  2.28% 62.47%      590ms  2.99%  runtime.findObject
```

### License
LGPL3