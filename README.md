![Logo](logo.png)

### Intro
g-fu is a pragmatic [Lisp](https://xkcd.com/297/) developed and embedded in Go. The initial [release](https://github.com/codr7/g-fu/tree/master/v1) implements an extensible, tree-walking interpreter with support for quasi-quotation and macros, lambdas, tail-call optimization, opt-/varargs, threads and channels; all weighing in at 2 kloc.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.11

Press Return twice to evaluate.

  (fun fib (n)
    (if (< n 2)
      n
      (+ (fib (- n 1)) (fib (- n 2)))))
```
```
  (fib 20))

6765
```

### Branching ([lib/cond.gf](https://github.com/codr7/g-fu/blob/master/v1/lib/cond.gf))

Every value has a boolean representation that may be retrieved using `bool`.

```
  (bool 42)

T

  (bool "")

F
```

Values may be combined using `or`/`and`. Unused values are not evaluated, and comparisons are performed using boolean representations while preserving the original values.

```
  (or 0 42)

42

  (or 0 F)
_

  (and '(1 2) '(3 4))

(3 4)

  (and '(1 2) F)
_
```

`if` may be used to branch on a condition.

```
  (if 42 'foo 'bar)

'foo
```

The else-branch is optional.

```

  (if "" 'foo)  
_
```

`switch` may be used to combine multiple branches.

```
  (switch
    (F 'foo)
    (T 'bar)
    (T 'baz))

'bar
```

### Iterating ([lib/iter.gf](https://github.com/codr7/g-fu/blob/master/v1/lib/iter.gf))
All loops support exiting with a result using `(break ...)` and skipping to the start of next iteration using `(continue)`.

The most fundamental loop is called `loop`, and that's exactly what it does until exited using `break` or external means such as `recall` and `fail`.
```
  (dump (loop (dump 'foo) (break 'bar) (dump 'baz)))

'foo
'bar
```

The `while`-loop keeps iterating until the specified condition turns false.

```
  (let (i 0)
    (while (< i 3)
      (dump (inc i))))

1
2
3
```

The `for`-loop accepts any iterable and an optional variable name, and runs one iteration for each value.

```
  (for 3 (dump 'foo))

'foo
'foo
'foo
```

```
  (for ('(foo bar baz) v) (dump v))

'foo
'bar
'baz
```

### Multitasking
Tasks are first class, preemptive green threads (or goroutines) that run in separate environments and interact with the outside world using channels. New tasks are started using `task` which optionally takes a task id and channel or buffer size argument and returns the new task. `wait` may be used to wait for task completion and get the results.

```
  (let _
    (task t1 () (dump 'foo) 'bar)
    (task t2 () (dump 'baz) 'qux)
    (dump (wait t1 t2)))

baz
foo
(bar qux)
```

The defining environment is cloned by default to prevent data races.

```
  (let (v 42)
    (dump (wait (task () (inc v))))
    (dump v))

43
42
```

#### Channels
Channels are optionally buffered, thread-safe pipes. `chan` may be used to create new channels, and `push`/`pop` to transfer values; `len` returns the current number of buffered values. Values sent over channels are cloned by default to prevent data races.

```
  (let (c (chan 1))
    (push c 42)
    (dump (len c))
    (dump (pop c)))

1
42
```

Unbuffered channels are useful for synchronizing tasks. The following example starts with the main task (which is unbuffered by default) `post`-ing itself to the newly started task `t`, which then replies `'foo` and finally returns `'bar`

```
  (let _
    (task t ()
      (post (fetch) 'foo)
      'bar)
      
    (post t (this-task))
    (dump (fetch))
    (dump (wait t)))

foo
bar
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