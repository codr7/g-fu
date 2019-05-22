## The DOOM Fire Kata

### Intro
Ever since I came across the DOOM fire [trick](https://fabiensanglard.net/doom_fire_psx/), I've been itching to work my way through it using console graphics for use as a kata to exercise new languages. This post describes how I would perform it in [g-fu](https://github.com/codr7/g-fu/tree/master/v1), a pragmatic Lisp embedded in Go.

![Fire](fire.gif)
[Source](https://github.com/codr7/g-fu/blob/master/v1/demo/fire.gf)

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
```

### Syntax
g-fu quasi-quotes using `'` and splices using `%`, `_` is used for missing values and `..` to splat sequences.

### The Idea
The idea is to model each particle of the fire as a value that decays along a reddish color scale while moving upwards. This is the reason for the white line at the bottom, that's where new particles are born. Add a touch of pseudo-chaos to make it interesting and that's pretty much it.

### Performance
While there's nothing seriously wrong with this implementation from what I can see, it's not going to win any performance prizes. [g-fu](https://github.com/codr7/g-fu/tree/master/v1) is still very young and I'm mostly working on correctness at this point. More mature languages with comparable features should be able to run this plenty faster.

Until next time,<br/>
c7