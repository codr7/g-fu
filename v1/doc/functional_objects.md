## Functional Objects

### Intro
This document describes the implementation of a minimal viable single-dispatch object system using closures in g-fu (https://github.com/codr7/g-fu), a pragmatic Lisp embedded in Go.

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
  (on-click)

  (resize (dx dy)
    (say "Button resize")
    (this 'Widget/resize dx dy))

  (on-click (f)
    (push on-click f))

  (click ()
    (for (on-click f) (f this))))
```
```
  (let (b (Button 'new 'width 100 'height 50))
    (say (b 'move 20 10))
    (say (b 'resize 100 0))
    (b 'on-click (fun (b) (say "Button click")))
    (b 'click))

20 10
Button resize
200 50
Button click
```

All of this and more may be loaded by evaluating `(load "doc/fos.gf")` from the release root.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.11

Press Return twice to evaluate.

  (load "doc/fos.gf")
_
```

### Syntax
g-fu quasi-quotes using `'` and splices using `%`, `_` is used for missing values and `..` to splat sequences.

### Dispatch
From one angle, a closure is essentially a single method object that uses its environment as storage. Adding a method argument and a `switch` extends the idea to support multiple methods. The `dispatch`-macro captures this pattern without assuming or dictating anything concerning object storage. `tr` is the standard tool for transforming sequences in g-fu; it takes an input, initial result and transformer.

```
(mac dispatch (defs..)
  (let args (new-sym) id (new-sym))
  
  '(fun (%args..)
     (let %id (head %args))
     
     (switch
       %(tr defs ()
            (fun (acc d)
              (let did (head d) imp (tail d))
              (push acc
                    (if (T? did)
                      '(T
                         (call (fun (%(head imp)..) %(tail imp)..)
                               %args..))
                      '((= %id '%did)
                         (call (fun (%(head imp)..) %(tail imp)..)
                               (splat (tail %args))))))))..
                            
       (T (fail (str "Unknown method: " %id))))))
```

The following exmple uses `let` to create a new environment containing a slot and `dispatch` to wrap it in a protocol. Calls to non-existing methods may be trapped by declaring a `T` method.

```
  (let (n 0 d (dispatch
                (inc ((delta 1)) (inc n delta))
                (T (args..) (say "Trapped: " args))))
    (say (d 'inc 42))
    (say (d 'inc))
    (say (d 'inc -1))
    (d 1 2 3))

42
43
42
Trapped: 1 2 3
```

Calls may be expanded to visually inspect the generated code.

```
  (expand 1 '(dispatch
               (inc ((delta 1)) (inc n delta))))

(fun (sym-143..)
  (let sym-144 (head sym-143))
  
  (switch
    ((= sym-144 'inc)
      (call (fun ((delta 1)) (inc n delta))
            (splat (tail sym-143))))
              
    (T (fail (str "Unknown method: " sym-144)))))
```

### This
`dispatch` is a useful tool in itself, but there comes a time when `this` needs to be accessed from the inside, to delegate etc. `let-this` expands to a new environment with `this` bound to the result of evaluating the last form in its body.

```
(mac let-this (vars body..)
  '(let (this _ %vars..)
     (set this %(pop body))
     %body..
     (fun (args..) (this args..))))
```
```
  (let-this ()
    (say this)
    42)

42
```

The following example creates a `this`-aware `dispatch`-er with a `patch`-method that may be used to hook into method dispatch and/or install new environments. In this case we're adding an interceptor that returns 42 when called without arguments and delegates anything else.

```
(let (s (let-this ()
          (dispatch
            (patch (new) (set this new)))))
  (s 'patch (fun (args..)
    (if args (s args..) 42)))
    
  (say (s)))

42
```

### Classification
Classes may be created using the `class`-macro. Classes are implemented as `this`-aware `dispatch`-ers, the constructor is just another method.

```
(mac class (id supers slots methods..)
  '(let %id
     (let-this ()
       (dispatch
         (id () '%id)
         (slots () '%slots)
         (methods () '%methods)
         (new (args..)
           (new-object (vec %supers..) '%slots '%methods args))))))
```

`eval` is used to create new objects to avoid requiring a central class registry and support lexically scoped class types. Super slots are prepended to the object's bindings, and super methods appended to the dispatch table. Super methods additionally support fully qualified names to allow delegation within overrides. Slot values passed to the constructor override init-forms.

```
(fun new-object (supers slots methods args)
  (eval '(let-this %(tr (push (super-slots supers) slots..) ()
                        (fun (acc x)
                          (if (= (type x) Vec)
                            (let (id (head x) v (pop-key args id))
                              (if (_? v) (push acc x..) (push acc id v)))
                            (push acc x (pop-key args x)))))
    %(and args (fail (str "Unused args: " args)))
    
    (dispatch
      %methods..
      %(super-methods supers)..))))
```

The task of collecting super slots makes a good match for [transducers](https://github.com/codr7/g-fu/blob/master/v1/lib/iter.gf). `t@` takes a reducing function as first argument which is used to instantiate the composed transformation pipeline, `tmap` followed by `tcat` in the following example.

```
(fun super-slots (supers)
  (tr supers () (t@ push (tmap (fun (s) (s 'slots))) tcat)))
```

Two dispatch entries are generated for each super method, one regular and one qualified with the super class name.

```
(fun super-methods (supers)
  (tr supers ()
      (fun (acc s)
        (tr (s 'methods) ()
            (fun (acc m)
              (push acc m '(%(sym (s 'id) '/ (head m)) %(tail m)..))))))))
```

And that's about it for now.<br/>

As you have probably guessed, the sky is the limit. I consider the functionality described here to be critical for a useful object system, and that includes convenient ways of hooking into the system to add additional layers of functionality.<br/>

Until next time,<br/>
c7