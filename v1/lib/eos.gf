(debug)

(load "../lib/abc.gf")
(load "../lib/iter.gf")

(fun new-object (slots args)
  '(let this (this-env)
        %(tr slots _
             (fun (acc x)
               (if (= (type x) Vec)
                 (let (id (head x) v (pop-key args id))
                   (if (_? v) (push acc x..) (push acc id v)))
                 (push acc x (pop-key args x)))))..))

(mac class (id supers slots methods..)  
  '(let %id (let (id '%id
                  this-class (this-env)
                  supers (vec %supers..)
                  slots (push (tr supers _
                                  (t@ push (tmap (fun (s) s/slots)) tcat))
                              '%slots..)
                  methods '%methods)

              (fun new (args..)
                (use _ do fail let this-class this-env)
                (eval (new-object slots args))
                this)
                
              this-class)))

(class Widget ()
  ((left 0) (top 0)
   (width (fail "Missing width")) (height (fail "Missing height"))))

(class Button (Widget)
  (on-click))

(dump Button/id)
(dump (Button/new 'width 100 'height 50))