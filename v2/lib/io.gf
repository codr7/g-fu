(fun say (args..)
  (for (args a)
    (print stdout a \n))
  (flush stdout))

(fun dump (args..)
  (for (args a)
    (print stderr a \n))
  (flush stderr))