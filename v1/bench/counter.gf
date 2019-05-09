(debug)
(load "../lib/all.gf")

(let _

(let super this-env
     Counter (fun (n)
               (use super inc dec)
               (fun c-inc () (inc n))
               (fun c-dec () (dec n))
               this-env))

(dump (bench 10 (for 1000
  (let c (Counter 0))
  (c/c-inc)
  (c/c-dec))))

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