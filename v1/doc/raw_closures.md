## Raw Closures

### Intro
A closure usually means a function with a derived environment. The environment is used to store data and closure/s form the api. So far so good, it's possible to build mostly anything using [TCO](http://wiki.c2.com/?TailCallOptimization)-gifted closures. [g-fu](https://github.com/codr7/g-fu/tree/master/v1) includes a basic single dispatch object [system](https://github.com/codr7/g-fu/blob/master/v1/doc/functional_objects.md) based on closures.

The thing that bugs me with functional closures is that they require squeezing all interaction through a functional/applicative pipe, which leads to reinventing unreachable features using suboptimal tools.

Peeling off the functional layer leaves the environment, which is promoted to first class.

```
```

Qualified ids allows reaching into external environments to access their slots/functions.

```
```