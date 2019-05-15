(debug)
(load "../lib/all.gf")

(let _

(let super this-env
     Counter (fun ((n 0))
               (fun inc ((d 1)) (super/inc n d))
               this-env))

(dump (bench 10 (for 1000
  (let c (Counter))
  (c/inc)
  (c/inc -1))))

)

(let _

(class Counter ()
  ((n 0))
  (inc ((d 1)) (inc n d)))

(dump (bench 10 (for 1000
  (let c (Counter 'new))
  (c 'inc)
  (c 'inc -1))))

)