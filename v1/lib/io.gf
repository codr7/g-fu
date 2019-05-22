(fun say (args..)
  (for (args a)
    (print stdout a LF))
  (flush stdout))

(fun dump (args..)
  (for (args a)
    (print stderr a LF))
  (flush stderr))