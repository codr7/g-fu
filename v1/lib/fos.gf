(let dispatch (mac (defs..)
  (let args (new-sym) id (new-sym))
  
  '(fun (%(args..))
     (let %id (head %args))
     
     (switch
       %(fold defs
              (fun (acc d)
                (let did (head d) imp (tail d))
                (push acc '(%(if (= did T) T '(= %id '%did))
                            ((fun (%(head imp)..) %(tail imp)..)
                             (splat (tail %args))))))
              _)..
       (T (fail "Unknown message"))))))

(let let-self (mac (vars body..)
  '(let (self _ %vars..)
     (set 'self %(pop body))
     %body..
     (fun (args..) (self args..)))))

(let class (mac (id supers slots methods..)
  '(let %id
     (let-self ()
       (dispatch
         (id () '%id)
         (slots () '%slots)
         (methods () '%methods)
         (new ()
           (let-self %(fold slots
                            (fun (acc x)
                              (if (== (type x) Vec) (push acc x..) (push acc x _)))
                            _)
             (dispatch %methods..))))))))
