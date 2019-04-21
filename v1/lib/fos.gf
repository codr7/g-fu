(let funs (mac (defs..)
  (let args (g-sym))
  
  '(fun (%(args..))
     (switch (head %args)
       %(fold defs
              (fun (acc d)
                (let imp (tail d))
                (push acc '((= '%(head d))
                            ((fun (%(head imp)..) %(tail imp)..)
                             (splat (tail %args))))))
              _)..
       (T (fail "Unknown message"))))))

(let lets (mac (vars body..)
  '(let (self _ %vars..)
     (let self %(pop body))
     %body..
     (fun (args..) (self args..)))))