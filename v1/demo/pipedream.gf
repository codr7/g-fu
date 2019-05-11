(debug)
(load "../lib/all.gf")

(let Port (let _
  (fun new ()
    (let this this-env io _
         elevation 0 sg 0 pressure 0
         init (fun () (if io (io/init))))
    this)

  this-env))

(fun connect (ports..)
  (if ports (do
    (let x (pop ports) y (pop ports))
    (set 'x/io y 'y/io x)
    (recall ports..))))

(mac define-node (id)
  '(let %id (let (node-type '%id)
     (fun new ((id (str '%id)))
       (let this this-env
            in (Port/new) out (Port/new))
       (set 'in/init (fun ()
                       (if in/io (set 'in/elevation in/io/elevation))
                       (set 'out/elevation in/elevation)
                       (out/init)))
       this)

     this-env)))

(define-node Pipe)

(define-node Valve)

(define-node Tank)

(let in-tank (Tank/new "In Tank")
     in-pipe (Pipe/new "In Pipe")
     valve (Valve/new)
     out-pipe (Pipe/new "Out Pipe")
     out-tank (Tank/new "Out Tank"))

(dump valve/id)

(connect in-tank/out in-pipe/in
         in-pipe/out valve/in
         valve/out out-pipe/in
         out-pipe/out out-tank/in)

(set 'in-tank/out/elevation 10)
(in-tank/out/init)
(dump out-tank/in/elevation)