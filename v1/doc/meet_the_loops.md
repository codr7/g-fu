## Meet the Loops

### Intro
This document describes the implementation of three fundamental loop constructs using closures and macros in g-fu (https://github.com/codr7/g-fu), a pragmatic Lisp embedded in Go.

All of this and more may be loaded by evaluating `(load "lib/all.gf")` from the release root.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.11

Press Return twice to evaluate.

  (load "lib/all.gf")
_
```

### Syntax
g-fu uses `%` as opposed to `,` for splicing, `_` in place of `nil`; and `..` to splat, which replaces `@`.

### while
First out is the `while`-loop, which keeps iterating until its condition turns false. It is implemented in terms of the more general `loop` macro which will follow shortly.

```
(mac while (cond body..)
  '(loop
     (if %cond _ (break))
     %body..))
```
```
  (let (i 0)
    (while (< i 3)
      (dump (inc i))))

1
2
3
```

### loop
The most fundamental loop is called `loop`. It supports skipping to the start of next iteration using `continue` and exiting with a result using `break`. Most of the work is performed by an anonymous, tail-recursive function; fresh argument symbols are created to avoid capturing the calling environment. g-fu supports explicit tail recursion using `recall`, the new call replaces the current one regardless of where it is triggered.

```
(mac loop (body..)
  (let done? (new-sym) result (new-sym))
  
  '(let _
     (mac break (args..) '(recall T %args..))
     (mac continue () '(recall))
     
     ((fun ((%done? F) %result..)
        (if %done?
          %result..
          (do %body.. (recall)))))))
```
```
  (dump (loop (dump 'foo) (break 'bar) (dump 'baz)))

'foo
'bar
```

### for
The `for`-loop accepts any iterable and an optional variable name, and runs one iteration for each value until the iterator returns `_`. Like `while`, it is based on `loop`.

```
(mac for (args body..)
  (let v? (= (type args) Vec)
       in (new-sym)
       out (if (and v? (> (len args) 1)) (pop args) (new-sym)))
       
  '(let (%in (iter %(if v? (pop args) args)))
     (loop
       (let %out (pop %in))
       (if (_? %out) (break))
       %body..)))
```
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

I feel like that's enough repetition for today.<br/>

Until next time,<br/>
c7