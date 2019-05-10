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
g-fu v1.11

Press Return twice to evaluate.

  (load "lib/all.gf")
_
```

### Syntax
g-fu uses `_` in place of `nil`, `..` to splat; and `/` as separator in qualified ids.

### The Environment
Peeling the functional layer off a closure leaves the (now first class) environment. `this-env` evaluates to the current environment. Note that used bindings,`this-env` in the following example, are automatically captured as usual.

```
  (let (foo 42) this-env)

(foo:42 this-env:(prim this-env))
```

Dealing directly with the environment allows composing data and code in a more flexible, performant and convenient form.

```
  (let (super this-env
        Counter (fun ((n 0))
                  (fun inc () (super/inc n))
                  (fun dec () (super/dec n))
                  this-env)
        c (Counter))
    (for 3 (c/inc))
    (c/dec))

2
```

Besides capturing used bindings, environments also support manually importing non-captured bindings with `use`.

```
(let (foo (let (bar 42)
            this-env)
      baz (let _
            (use foo bar)
            this-env))
  baz/bar)

42
```

Failed lookups may be trapped by defining `resolve`. The following example implements a basic proxy that forwards all lookups to the specified delegate.

```
(fun proxy (d)
  (fun resolve (id)
    (d/val id))

  this-env)

```
```
  (let (p (proxy (let (foo 42) this-env)))
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
  
    this-env)

  this-env))
```

Buttons embed a Widget, delegate `move` and override `resize` to enforce a max size. They also support registering `on-click`-handlers.

```
(let Button (let _
  (fun new (args..)
    (let this (this-env)
         w (Widget/new args..)
         click-event _)
         
    (use w move)

    (fun click ()
      (for (click-event f) (f this)))
      
    (fun on-click (f)
      (push click-event f))

    (fun resize (dx dy)
      (w/resize (min (+ w/width dx) (- 200 w/width))
                (min (+ w/height dy) (- 100 w/height))))
    
    this)

  this-env))
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