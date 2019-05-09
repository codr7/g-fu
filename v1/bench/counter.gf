(debug)
(load "../lib/all.gf")

(let _

(let super this-env
     Counter (fun ((n 0))
               (fun inc () (super/inc n))
               (fun dec () (super/dec n))
               this-env))

(dump (bench 10 (for 1000
  (let c (Counter))
  (c/inc)
  (c/dec))))

)

(let _

(class Counter ()
  ((n 0))
  (inc () (inc n))
  (dec () (dec n)))

(dump (bench 10 (for 1000
  (let c (Counter 'new))
  (c 'inc)
  (c 'dec))))

)