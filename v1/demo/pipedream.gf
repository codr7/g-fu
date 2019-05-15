(debug)
(load "../lib/all.gf")

(fun port (n)
  (let this this-env
       y .0 dy .0
       pressure .0 density 1.
       io _
       default (let _
                 (fun init (prev)
                   (set 'y (+ prev/y dy))
                   (if (= n prev) _ (n/init))
                   (if (and io (not (= io prev))) (io/init this)))

                 (fun pair (p)
                   (set 'io p 'p/io this))
    
                 (fun sim (prev)
                   (dump 'port-sim)
                   (set 'pressure prev/pressure)
                   (if (= n prev) _ (n/sim))
                   (if (and io (not (= io prev))) (io/sim this)))
                   
                 this-env))
  (use default init pair sim)
  this)

(mac node (id (vars ()) body..)
  '(fun %id ((id (str '%id)))
     (let this this-env
          y .0 dy .0
          pressure .0 density 1.
          in (port this) out (port this)
          %vars..
          default (let _
                    (fun init ()
                      (set 'y (+ in/y dy))
                      (out/init this))

                    (fun sim ()
                      (in/sim this)
                      (out/sim this))

                    this-env))
     (use default init sim)
     %body..
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

(fun y->pressure (y)
  (* y 1.422))

(node Pipe
  (diameter .0 length .0 flow .0)

  (fun get-pressure ()
    (- in/pressure
       (* (* (/ density 2.) (/ (* flow) diameter)) length)))

  (fun sim ()
    (dump 'pipe-sim)
    (set 'pressure (get-pressure))
    (out/sim this)))

(node Tank
  (volume .0 radius .0)

  (fun get-pressure ((dy .0))
    (y->pressure (- (/ volume (* PI (* radius))) dy)))

  (fun sim ()
    (dump 'tank-sim)
    (set 'pressure (get-pressure out/dy))
    (out/sim this)))
    
(node Valve)

(let-node t1 (Tank) p1 (Pipe) v (Valve) p2 (Pipe) t2 (Tank))
(__ (chain t1 p1 v p2 t2))
(chain t1 p1)

(set 't1/radius 10. 't1/volume 10000. 't1/dy 10.
     'p1/diameter .1 'p1/length 10.)

(t1/init)
(t1/sim)
(dump p1/in/pressure p1/out/pressure)