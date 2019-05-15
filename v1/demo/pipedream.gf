(debug)
(load "../lib/all.gf")

(fun Port (n)
  (let this this-env
       y .0 dy .0
       io _
       default (let _
                 (fun init (prev)
                   (set 'y (+ prev/y dy))
                   (if (= n prev) _ (n/init))
                   (if (and io (not (= io prev))) (io/init this)))

                 (fun pair (p)
                   (set 'io p 'p/io this))
    
                 this-env))
  (use default init pair)
  this)

(mac define-node (id)
  '(fun %id ((id (str '%id)))
     (let this this-env
          y .0 dy .0
          sg .0 pressure .0
          in (Port this) out (Port this)
          default (let _
                    (fun init ()
                      (set 'y (+ in/y dy))
                      (out/init this))

                    this-env))
     (use default init)
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
(set 't1/dy 10.)
(t1/init)
(dump t2/y)
