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

(let in-tank (Tank "In Tank")
     in-pipe (Pipe "In Pipe")
     valve (Valve)
     out-pipe (Pipe "Out Pipe")
     out-tank (Tank "Out Tank"))

(dump valve/id)

(connect in-tank/out in-pipe/in
         in-pipe/out valve/in
         valve/out out-pipe/in
         out-pipe/out out-tank/in)

(set 'in-tank/out/elevation 10)
(in-tank/out/init)
(dump out-tank/in/elevation)