(debug)
(load "../lib/all.gf")

(fun Port ()
  (let this this-env
       io _
       elevation .0 sg .0 pressure .0
       default (let _
                 (fun init ()
                   (if io (io/init)))

                 (fun pair (p)
                   (set 'io p 'p/io this))
    
                 this-env))
  (use default init pair)
  this)

(mac define-node (id)
  '(fun %id ((id (str '%id)))
     (let this this-env
          in (Port) out (Port))
          
     (fun in/init ()
       (if in/io (set 'in/elevation in/io/elevation))
       (set 'out/elevation in/elevation)
       (out/init))
       
     this))

(mac let-node (args..)
  (let a1 (head args))
  
  (fun push-args (in (acc ()))
    (if in (do
      (let v (pop in) k (pop in))
      (push acc k '(%(head v) '%k %(tail v)..))
      (recall in acc))
      acc))
    
  (if (Vec a1)
    '(let %(push-args a1) (tail args)..)
    '(let %(push-args args)..)))

(fun chain (ns..)
  (tr (tail ns) (head ns)
      (fun (x y)
        (x/out/pair y/in)
        y)))

(define-node Pipe)
(define-node Tank)
(define-node Valve)

(let-node t1 (Tank) p1 (Pipe) v (Valve) p2 (Pipe) t2 (Tank))
(chain t1 p1 v p2 t2)

(set 't1/out/elevation 10.)
(t1/out/init)
(dump t2/in/elevation)