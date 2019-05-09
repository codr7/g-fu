## Raw Closures

### Intro
The word Closure usually means a function with a captured environment. So far so good, it's possible to build mostly anything using [TCO](http://wiki.c2.com/?TailCallOptimization)-gifted closures. [g-fu](https://github.com/codr7/g-fu/tree/master/v1) includes a basic single dispatch object [system](https://github.com/codr7/g-fu/blob/master/v1/doc/functional_objects.md) based on closures.

The thing that been increasingly bugging me is being required to squeeze all interaction with the underlying environment through a function, which is both slow and inconvenient. 

Peeling off the functional layer leaves the (now first class) environment. Note that used bindings,`this-env` in the following example, are captured as usual.

```
  (let _ this-env)

(this-env:(prim this-env ()))
```

Qualified ids allow reaching into external environments to access their bindings.

```
```

Combining these ideas allows combining data and code in a more flexible, performant and convenient form.

```
(let (super this-env
      Counter (fun ((n 0))
                (fun inc () (super/inc n))
                (fun dec () (super/dec n))
                this-env)
      c (Counter))
  (c/inc)
  (c/inc)
  (c/inc)
  (c/dec))

2
```

(use)

```
```

GUI

```
```