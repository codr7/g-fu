## Raw Closures

### Intro
Closures are functions with captured environments. So far so good, it's possible to build mostly anything using [TCO](http://wiki.c2.com/?TailCallOptimization)-gifted closures. [g-fu](https://github.com/codr7/g-fu/tree/master/v1) includes a single dispatch object [system](https://github.com/codr7/g-fu/blob/master/v1/doc/functional_objects.md) based on nothing but closures.

The thing that's been increasingly bugging me is being required to squeeze all interaction with the underlying environment through a function, which is both slow and inconvenient in all but the most trivial cases.

### Setup
If you feel like playing around with the code, the following shell spell should put you in a REPL that supports all examples.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.12

Press Return twice to evaluate.

  (load "lib/all.gf")
_
```

### Syntax
g-fu quasi-quotes using `'` and splices using `%`, `_` is used for missing values and `..` to splat sequences.

### The Environment
Peeling the functional layer off a closure leaves the (now first class) environment. `Env/this` evaluates to the current environment. Note that used bindings,`Env/this` in the following example, are automatically captured as usual.

```
  (let (foo 42) Env/this)

(foo:42 Env/this:(prim Env/this))
```

Dealing directly with the environment allows composing data and code in a more flexible, performant and convenient form.

```
  (let (super Env/this
        Counter (fun ((n 0))
                  (fun inc ((d 1)) (super/inc n d))
                  Env/this)
        c (Counter))
    (for 3 (c/inc))
    (c/inc -1))

2
```

Non-captured bindings may be imported manually.

```
  (let _
    (env foo (bar 42))
    (env baz _ (use foo bar))
    baz/bar)

42
```

Failed lookups may be trapped by defining `resolve`. The following example implements a basic proxy that forwards all lookups to the specified delegate.

```
(fun proxy (d)
  (fun resolve (key)
    (d/val key))

  Env/this)
```
```
  (let (p (proxy (let (foo 42)
                   (use _ val)
                   Env/this)))
    p/foo)

42
```

We will end this environmental adventure in the land of graphical user interfaces. A widget knows its position and dimensions, and supports moving and resizing. The type, or class; is just another environment containing the constructor and potentially more static methods and/or fields.

```
(let Widget (let _
  (fun new (args..)
    (let left 0 top 0
         width (or (pop-key args 'width) (fail "Missing width"))
         height (or (pop-key args 'height) (fail "Missing height")))

    (fun move (dx dy)
      (vec (inc left dx)
           (inc top dy)))

    (fun resize (dx dy)
      (vec (inc width dx)
           (inc height dy)))
  
    Env/this)

  Env/this))
```

Buttons embed a Widget, delegate `move` and override `resize` to enforce a max size. They also support registering `on-click`-handlers.

```
(let Button (let _
  (fun new (args..)
    (let w (Widget/new args..)
         click-event ())
         
    (use w move)

    (fun click ()
      (for (click-event f) (f Env/this)))
      
    (fun on-click (f)
      (push click-event f))

    (fun resize (dx dy)
      (w/resize (min (+ w/width dx) (- 200 w/width))
                (min (+ w/height dy) (- 100 w/height))))
    
    Env/this)

  Env/this))
```
```
  (let (b (Button/new 'width 100 'height 50))
    (say (b/move 10 10))
    (say (b/resize 400 200))
    (b/on-click (fun (b) (say "Button clicked!")))
    (b/click))

10 10
200 100
Button clicked!
```

Until next time,<br/>
c7