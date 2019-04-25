(let dispatch (mac (defs..)
  (let args (new-sym) id (new-sym))
  
  '(fun (%args..)
     (let %id (head %args))
     
     (switch
       %(fold defs
              (fun (acc d)
                (let did (head d) imp (tail d))
                (push acc
                      (if (= did T)
                        '(T
                           ((fun (%(head imp)..) %(tail imp)..)
                             %args..))
                        '((= %id '%did)
                           ((fun (%(head imp)..) %(tail imp)..)
                            (splat (tail %args))))))))..
       (T (fail (str "Unknown message: " %id)))))))

(let let-self (mac (vars body..)
  '(let (self _ %vars..)
     (set 'self %(pop body))
     %body..
     (fun (args..) (self args..)))))

(let super-slots (fun (supers)
  (fold supers (fun (acc s) (push acc (s 'slots)..)))))

(let super-methods (fun (supers)
  (fold supers
        (fun (acc s)
          (fold (s 'methods)
                (fun (acc m)
                  (push acc m '(%(sym (s 'id) '/ (head m)) %(tail m)..))))))))

(let new-object (fun (supers slots methods args)
  (eval '(let-self %(fold (append (super-slots supers) slots..)
                          (fun (acc x)
                            (if (= (type x) Vec)
                              (let (id (head x) v (find-key args id))
                                (if (= v _) (push acc x..) (push acc id v)))
                              (push acc x (find-key args x)))))
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
