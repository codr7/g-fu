## The DOOM Fire Kata

### Intro
Ever since I came across the DOOM fire [trick](https://fabiensanglard.net/doom_fire_psx/), I've been itching to work my way through it using console graphics for use as a kata to exercise programming languages. This post describes how I would do it in [g-fu](https://github.com/codr7/g-fu/tree/master/v1), a pragmatic Lisp embedded in Go.

![Fire](fire.gif)

### Setup
If you feel like lighting your own fire, the following shell spell will take you where you need to go.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.15

Press Return twice to evaluate.

  (load "demo/fire.gf")
_
```

### Syntax
g-fu quasi-quotes using `'` and splices using `%`, `_` is used for missing values and `..` to splat sequences.

Until next time,<br/>
c7