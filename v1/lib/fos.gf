(let dispatch (mac (defs..)
  (let args (new-sym) id (new-sym))
  
  '(fun (%args..)
     (let %id (head %args))
     
     (switch
       %(fold defs _
              (fun (acc d)
                (let did (head d) imp (tail d))
                (push acc
                      (if (T? did)
                        '(T
                           ((fun (%(head imp)..) %(tail imp)..)
                             %args..))
                        '((= %id '%did)
                           ((fun (%(head imp)..) %(tail imp)..)
                            (splat (tail %args))))))))..
       (T (fail (str "Unknown method: " %id)))))))

(let let-self (mac (vars body..)
  '(let (self _ %vars..)
     (set 'self %(pop body))
     %body..
     (fun (args..) (self args..)))))

(let super-slots (fun (supers)
  (fold supers _ (@ push (map-r (fun (s) (s 'slots))) cat-r))))

(let super-methods (fun (supers)
  (fold supers _
        (fun (acc s)
          (fold (s 'methods) _
                (fun (acc m)
                  (push acc m '(%(sym (s 'id) '/ (head m)) %(tail m)..))))))))

(let new-object (fun (supers slots methods args)
  (eval '(let-self %(fold (push (super-slots supers) slots..) _
                          (fun (acc x)
                            (if (= (type x) Vec)
                              (let (id (head x) v (pop-key args id))
                                (if (_? v) (push acc x..) (push acc id v)))
                              (push acc x (pop-key args x)))))
    %(and args (fail (str "Unused args: " args)))
    (dispatch
      %methods..
      %(super-methods supers)..)))))

(let class (mac (id supers slots methods..)
  '(let %id
     (let-self ()
       (dispatch
         (id () '%id)
         (slots () '%slots)
         (methods () '%methods)
         (new (args..)
           (new-object (vec %supers..) '%slots '%methods args)))))))
