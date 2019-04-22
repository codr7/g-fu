(let fo-fun (mac (defs..)
  (let args (new-sym))
  
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

(let fo-let (mac (vars body..)
  '(let (self _ %vars..)
     (set self %(pop body))
     %body..
     (fun (args..) (self args..)))))