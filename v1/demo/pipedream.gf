(debug)
(load "../lib/all.gf")

(fun Port ()
  (let this this-env io _
       elevation 0 sg 0 pressure 0
       init (fun () (if io (io/init))))
    this)

(fun connect (ports..)
  (if ports (do
    (let x (pop ports) y (pop ports))
    (set 'x/io y 'y/io x)
    (recall ports..))))

(mac define-node (id)
  '(fun %id ((id (str '%id)))
     (let this this-env
          in (Port) out (Port))
          
     (set 'in/init (fun ()
                     (if in/io (set 'in/elevation in/io/elevation))
                     (set 'out/elevation in/elevation)
                     (out/init)))
     this))

(define-node Pipe)

(define-node Valve)

(define-node Tank)

(let t1 (Tank 't1)
     p1 (Pipe 'p1)
     v (Valve 'v)
     p2 (Pipe 'p2)
     t2 (Tank 't2))

(connect t1/out p1/in
         p1/out v/in
         v/out p2/in
         p2/out t2/in)

(set 't1/out/elevation 10)
(t1/out/init)
(dump t2/in/elevation)