## Functional Objects

### Intro
This document describes the implementation of a minimal viable single-dispatch object system using closures and macros in g-fu (https://github.com/codr7/g-fu), a pragmatic Lisp embedded in Go.

```
  (class Widget ()
    ((left 0) (top 0)
     (width (fail "Missing width")) (height (fail "Missing height")))
  
    (move (dx dy)
      (vec (inc left dx)
           (inc top dy)))

    (resize (dx dy)
      (vec (inc width dx)
           (inc height dy))))

  (class Button (Widget)
    ((label "") on-click)

    (resize (dx dy)
      (say "Button resize")
      (self 'Widget/resize dx dy))

    (on-click (f)
      (push on-click f))

    (click ()
      (for (on-click f) (f self))))

  (let (b (Button 'new 'width 100 'height 50 'label "Click me"))
    (say (b 'move 20 10))
    (say (b 'resize 100 0))
    (b 'on-click (fun (b) (say "Button click")))
    (b 'click))

20 10
Button resize
200 50
Button click
```

All of this and more may be loaded by evaluating `(load "lib/all.gf")` from the release root.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.10

Press Return twice to evaluate.

  (load "lib/all.gf")
_
```

### Syntax
g-fu uses `%` as opposed to `,` for splicing, `_` in place of `nil`; and `..` to splat, which replaces `@`.

### Dispatch
From one angle, a closure is essentially a single method object that uses its environment as storage. Adding a method argument and a `switch` extends the idea to support multiple methods. The `dispatch`-macro captures this pattern without assuming anything about object storage.

```
  (let dispatch (mac (defs..)
    (let args (new-sym) id (new-sym))
  
    '(fun (%args..)
       (let %id (head %args))
     
       (switch
         %(fold defs _
                (fun (acc d)
                  (let did (head d) imp (tail d))
                  (push acc
                        (if (= did T)
                          '(T
                             ((fun (%(head imp)..) %(tail imp)..)
                               %args..))
                          '((= %id '%did)
                             ((fun (%(head imp)..) %(tail imp)..)
                              (splat (tail %args))))))))..
         (T (fail (str "Unknown method: " %id)))))))
```

The following exmple uses `let` to create a new environment containing a slot and `dispatch` to wrap it in a protocol. Calls to non-existing methods may be trapped by declaring a `T` method.

```
  (let (n 0 d (dispatch
                (inc ((delta 1)) (inc n delta))
                (dec ((delta 1)) (dec n delta))
                (T (args..) (say "Trapped: " args))))
    (say (d 'inc 42))
    (say (d 'inc))
    (say (d 'dec))
    (d 1 2 3))

42
43
42
Trapped: 1 2 3
```

Calls may be expanded to visually inspect the generated code.

```
  (expand 1 '(dispatch
               (inc ((delta 1)) (inc n delta))
               (dec ((delta 1)) (dec n delta))))

(fun (sym-143..)
  (let sym-144 (head sym-143))
    (switch
      ((= sym-144 'inc)
        ((fun ((delta 1)) (inc n delta)) (splat (tail sym-143))))
      ((= sym-144 'dec)
        ((fun ((delta 1)) (dec n delta)) (splat (tail sym-143))))
      (T (fail (str "Unknown method: " sym-144)))))
```

### Self
`dispatch` is a useful tool in itself, but there comes a time when `self` needs to be accessed from the inside, to delegate etc. `let-self` expands to a new environment with `self` bound to the result of evaluating the last form in its body.

```
(let let-self (mac (vars body..)
  '(let (self _ %vars..)
     (set 'self %(pop body))
     %body..
     (fun (args..) (self args..)))))
```
```
  (let-self ()
    (say self)
    42)

42
```

The following example creates a self-aware `dispatch` with a `patch`-method that may be used to hook into method dispatch and/or install a different environment. In this case we're installing an intercepting function that returns 42 when called without arguments and delegates anything else to the original dispatch table.

```
(let (s (let-self ()
          (dispatch
            (patch (new) (set 'self new)))))
  (s 'patch (fun (args..)
    (if args (s args..) 42)))
    
  (say (s)))

42
```

### Classification
Classes, or object factories; may be created using the `class`-macro. Classes are implemented as self-aware dispatchers, the constructor is just another method.

```
(let class (mac (id supers slots methods..)
  '(let %id
     (let-self ()
       (dispatch
         (id () '%id)
         (slots () '%slots)
         (methods () '%methods)
         (new (args..)
           (new-object (vec %supers..) '%slots '%methods args)))))))
```

g-fu uses `eval` for creating new objects without requiring a central class registry and to enable eventually supporting lexically scoped class types. Super slots are prepended to the object's bindings, and super methods appended to the dispatch table. Super methods additionally support fully qualified names to allow delegation within overrides. Slot values passed to the constructor override init-forms.

```
(let new-object (fun (supers slots methods args)
  (eval '(let-self %(fold (push (super-slots supers) slots..) _
                          (fun (acc x)
                            (if (= (type x) Vec)
                              (let (id (head x) v (pop-key args id))
                                (if (_? v) (push acc x..) (push acc id v)))
                              (push acc x (pop-key args x)))))
    %(and args (fail (str "Unused args: " args)))
    
    (dispatch
      %methods..
      %(super-methods supers)..)))))
```

The task of collecting super slots makes a good match for [transducers](https://github.com/codr7/g-fu/blob/master/v1/lib/iter.gf). `@` takes a reducing function as first argument and propagates it through the specified transformation pipeline, `map` followed by `cat` in the following example.

```
(let super-slots (fun (supers)
  (fold supers _ (@ push (map (fun (s) (s 'slots))) cat))))
```

Two dispatch entries are generated for each super method, one regular and one qualified with the super class name.

```
(let super-methods (fun (supers)
  (fold supers _
        (fun (acc s)
          (fold (s 'methods) _
                (fun (acc m)
                  (push acc m '(%(sym (s 'id) '/ (head m)) %(tail m)..))))))))
```

And that's about it for now.<br/>

As you have probably guessed, the sky is the limit. I consider the functionality described here to be critical for a useful object system, and that includes convenient ways of hooking into the system to add additional layers of functionality.<br/>

Until next time,<br/>
c7