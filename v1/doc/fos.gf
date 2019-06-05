(load "../lib/all.gf")

(mac dispatch (defs..)
  '(fun (%$args..)
     (let %$id (head %$args))
     
     (switch
       %(tr defs ()
            (fun (acc d)
              (let did (head d) imp (tail d))
              (push acc
                    (if (T? did)
                      '(T
                         (call (fun (%(head imp)..) %(tail imp)..)
                               %$args..))
                      '((= %$id '%did)
                         (call (fun (%(head imp)..) %(tail imp)..)
                               (splat (tail %$args))))))))..
                               
       (T (fail (str "Unknown method: " %$id))))))

(mac let-this (vars body..)
  '(let (this _ %vars..)
     (set this %(pop body))
     %body..
     (fun (args..) (this args..))))

(fun super-slots (supers)
  (tr supers () (t@ push (tmap (fun (s) (s 'slots))) tcat)))

(fun super-methods (supers)
  (tr supers ()
      (fun (acc s)
        (tr (s 'methods) ()
            (fun (acc m)
              (push acc m '(%(sym (s 'id) '/ (head m)) %(tail m)..)))))))

(fun new-object (supers slots methods args)
  (eval '(let-this %(tr (push (super-slots supers) slots..) ()
                        (fun (acc x)
                          (if (= (type x) Vec)
                            (let (id (head x) v (pop-key args id))
                              (if (_? v) (push acc x..) (push acc id v)))
                            (push acc x (pop-key args x)))))
    %(and args (fail (str "Unused args: " args)))
    (dispatch
      %methods..
      %(super-methods supers)..))))

(mac class (id supers slots methods..)
  '(let %id
     (let-this ()
       (dispatch
         (id () '%id)
         (slots () '%slots)
         (methods () '%methods)
         (new (args..)
           (new-object (vec %supers..) '%slots '%methods args))))))