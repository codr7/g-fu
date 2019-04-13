![Logo](../alien.png)

### Intro
This is the first in a planned series of posts documenting the implementation of [g-fu](https://github.com/codr7/g-fu), a pragmatic [Lisp](https://xkcd.com/297/) developed and embedded in Go. There are many Lisps in Go; but I find most too naive or dogmatic to be of practical use. g-fu is still in its infancy; and the current implementation leaves a lot to wish for, not least in the performance department. But what it does offer is an extensible tree-walking interpreter with support for quasi-quotation and macros, lambdas, bindings, tail-call optimization, opt-/varargs and eval; all weighing in below 2 kloc.

### Next
The main focus for next iteration is adding a compilation stage and an additional, more concrete internal code representation; which should improve performance considerably.

That's all for now, happy Lisping!
And please don't hesitate to ask questions and/or suggest improvements.

/c7