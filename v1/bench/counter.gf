(debug)
(load "../lib/all.gf")

(let _

(let super Env/this
     Counter (fun ((n 0))
               (fun inc ((d 1)) (super/inc n d))
               Env/this))

(dump (bench 10 (for 10000
  (let c (Counter))
  (c/inc)
  (c/inc -1))))

)