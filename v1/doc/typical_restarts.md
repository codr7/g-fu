## Typical Restarts

### Intro
Error handling has become a polarizing issue lately, arguments over merits of exceptions vs. manual propagation look more and more like dynamic vs. static types or Emacs vs. Vim every day. Take one step back and it looks like implicit vs. explicit, relaxed vs. disciplined or CISC vs. RISC. Different flavors of the same lovely ice cream.

The Issue with a capital I when it comes to error handling from my experience is that the code that knows what to do often is located several stack frames upstream from the crime scene.

Exceptions and manual propagation are simply different ways of passing enough information about the error upstream until it reaches a layer that knows what to do. The same information is then often passed down again in a separate call to the level where the error originated. I'm sure most would agree that the process feels a tiny bit more complicated than it needs to be.

Restarts allow passing information upstream without having to deal with unwinding stacks or manual propagation. Common Lisp's [condition system](http://www.gigamonkeys.com/book/beyond-exception-handling-conditions-and-restarts.html) is the only implementation I am aware of, though it muddies the water somewhat by throwing exceptions into the mix.

The idea is that code may provide options for dealing with errors for upstream handlers to choose from. Once a choice has been made, execution typically continues at the level where the error originated.

### Setup
If you feel like coding along, the following shell spell will take you where you need to go.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.20

Press Return twice to evaluate.

  (load "lib/all.gf")
  
```

### Throwing
Any value may be thrown up the call stack. When no other options remain for dealing with a `throw`, the system will enter a break loop which allow interactively invoking available restarts. `abort` and `retry` are always provided.

```
  (throw 42)

Break: 42
0 abort
1 retry

Choose 0-1: 1

Break: 42
0 abort
1 retry

Choose 0-1: 0

Abort
```

`fail` provides a convenient shortcut for throwing errors.

```
  (fail "Going down")
  
Break: Error: Going down
0 abort
1 retry

Choose 0-1:
```

### Trying
`try` may be used to limit the scope for `retry` and supports defining custom restarts. Restarts may declare any number of arguments. Execution typically continues after the `try`, returning the last value, once the restart exits.

```
  (try ((foo (x)
          (say (str "foo " x))
          'bar))
    (fail "Going down"))

Break: Error: Going down
0 abort
1 retry
2 foo x

Choose 0-2: 2 42

foo 42
bar
```

`restart` may be used to look up restarts in the current call stack.

```
  (try ((foo (x) (+ x 35)))
    (try _
      (call (restart 'foo) 7)))

42
```

Restarts live in a separate namespace, but allow shadowing just like regular bindings.

```
  (try ((foo () 'bar))
    (try ((foo () 'baz))
      (call (restart 'foo))))
  
baz
```

### Opening Files
A more realistic scenario may be triggered by loading a nonexistent file, which includes a restart for using a different filename.

test.gf
```
42
```

```
(load "not.found")

Break: Error: Failed loading file: "not.found"
open not.found: no such file or directory
0 abort
1 retry
2 use-filename new

Choose 0-2: 2 "test.gf"

42
```

If you still can't see the point, imagine an expensive computation before the `load` that is used after.

```
  (do
    (say "Expensive computation")
    (load "not.found")
    (say "Use result of computation"))      
  
Expensive computation

Break: Error: Failed loading file: "not.found"
open not.found: no such file or directory
0 abort
1 retry
2 use-filename new

Choose 0-2: 2 "test.gf"

Use result of computation
```

### Catching
The missing piece of the puzzle is a way to catch errors and invoke restarts programatically, which is where `catch` comes into the picture. The following example catches a symbol lookup error and provides a new value.

Let's start with a break loop to see what options are available.

```
  foo
  
Break: Error: Unknown: foo
0 abort
1 retry
2 use-key new
3 use-val new
```

Error handlers are expected to return a restart curried with any required arguments or `_` to enter a break loop.

```
  (catch (((EUnknown _) (restart 'use-val 42)))
    foo)

42
```

Until next time,<br/>
c7