(debug)
(load "../lib/all.gf")

(fun port (n)
  (let this this-env
       y .0 dy .0
       pressure .0 density 1.
       io _
       default (let _
                 (fun init (prev)
                   (set y (+ prev/y dy))
                   (if (= n prev) _ (n/init))
                   (if (and io (not (= io prev))) (io/init this)))

                 (fun pair (p)
                   (set io p p/io this))
    
                 (fun run (prev)
                   (dump 'port-run)
                   (set pressure prev/pressure)
                   (if (= n prev) _ (n/run))
                   (if (and io (not (= io prev))) (io/run this)))
                   
                 this-env))
  (use default init pair run)
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
                      (set y (+ in/y dy))
                      (out/init this))

                    (fun run ()
                      (in/run this)
                      (out/run this))

                    this-env))
     (use default init run)
     %body..
     this))

(mac sim (id args body..)
  (fun push-args (in (acc ()))
    (if in (do
      (let v (pop in) k (pop in))
      (push acc k '(%(head v) '%k %(tail v)..))
      (recall in acc))
      acc))
  
  '(let %id (let (this-sim this-env
                  %(push-args args)..)
              %body..
              this-sim)))

(fun chain (ns..)
  (tr (tail ns) (head ns)
      (fun (x y)
        (x/out/pair y/in)
        y)))

(fun height->psi (y)
  (* y 1.422))

(node Pipe
  (diameter .0 length .0 flow .0)

  (__ "Use Darcy–Weisbach equation solved for flow rate.")

  (fun run ()
    (dump 'pipe-run)
    (out/run this)))

(node Tank
  (volume .0 radius .0)

  (fun get-pressure ((dy .0))
    (- (/ volume (* PI (* radius))) dy))

  (fun run ()
    (dump 'tank-run)
    (set pressure (get-pressure out/dy))
    (out/run this)))
    
(node Valve)

(sim s (t1 (Tank) p1 (Pipe) v (Valve) p2 (Pipe) t2 (Tank))
  (chain t1 p1)

  (set t1/radius 10. t1/volume 10000.
       p1/diameter .1 p1/length 10.)

  (t1/init)
  (t1/run)
  (dump p1/in/pressure p1/out/pressure))