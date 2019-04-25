## Functional Objects

```
  (class Widget ()
    ((left 0) (top 0) width height)
  
    (move (dx dy)
      (inc left dx)
      (inc top dy))

    (resize (dx dy)
      (inc width dx)
      (inc height dy)))

  (class Button (Widget)
    ((label "") on-click)

    (on-click (f)
      (push on-click f))

    (click ()
      (fold on-click (fun (acc f) (f self)))))

  (let (b (Button 'new 'width 100 'height 50 'label "Click me!"))
    (b 'on-click (fun (b) (dump "Button clicked!")))
    (b 'click))

"Button clicked!"
```

This document contains a recipe for constructing a minimal viable object system using nothing but closures and macros. The code in the examples is taken from g-fu (https://github.com/codr7/g-fu), a Lisp-dialect implemented in Go. A full implementation of these ideas may be loaded by evaluating `(load "lib/all.gf")` from the release root.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.9

Press Return twice to evaluate.

  (load "lib/all.gf")
_
```

### Syntax
g-fu uses `%` as opposed to `,` for splicing, `_` in place of `nil`; and `..` to splat, which replaces `@`. With that out of the way, let's have a look at the implementation.

### Dispatching Methods
From one angle, a closure is essentially a single method object that uses its environment as storage. Adding a method argument and a `switch` extends the idea to support multiple methods. The `dispatch`-macro captures this pattern without assuming anything about object storage.

```
(let dispatch (mac (defs..)
  (let args (new-sym) id (new-sym))
  
  '(fun (%args..)
     (let %id (head %args))
     
     (switch
       %(fold defs
              (fun (acc d)
                (let did (head d) imp (tail d))
                (push acc '(%(if (= did T) T '(= %id '%did))
                            ((fun (%(head imp)..) %(tail imp)..)
                             (splat (tail %args))))))
              _)..
       (T (fail "Unknown message"))))))
```

The following exmple uses `let` to create an environment containing a slot and `dispatch` to wrap it in a simple protocol.

```
(let (n 0 d (dispatch
              (inc ((delta 1)) (inc n delta))
              (dec ((delta 1)) (dec n delta))))
  (dump (d 'inc 42) 42)
  (dump (d 'inc) 43)
  (dump (d 'dec) 42))
```

Expanding The call allows visually inspecting the generated code.

```
(dump (expand 1 '(dispatch
                  (inc ((delta 1)) (inc n delta))
                  (dec ((delta 1)) (dec n delta)))))
```
```
(fun (sym-135..)
  (let sym-136 (head sym-135))
  (switch
    ((= sym-136 'inc)
      ((fun ((delta 1)) (inc n delta)) (splat (tail sym-135))))
    ((= sym-136 'dec)
      ((fun ((delta 1)) (dec n delta)) (splat (tail sym-135))))
    (T (fail "Unknown message"))))
```

### Talking to Self
`dispatch` is a very useful tool, but there comes a time when you need access to `self` from from the inside, to delegate method calls etc. `let-self` expands to an environment with `self` bound to the result of evaluating its body.

```
(let let-self (mac (vars body..)
  '(let (self _ %vars..)
     (set 'self %(pop body))
     %body..
     (fun (args..) (self args..)))))
```
```
(let-self ()
  (dump self)
  42)

42
```

The following example creates a self-aware dispatcher with a `patch`-method that may be used to intercept method calls, in this case returning `42` when called without arguments.

```
(let (s (let-self ()
           (dispatch
             (patch (new) (set 'self new)))))
  (s 'patch (fun (args..)
    (if args (s args..) 42)))
    
  (dump (s)))

42
```