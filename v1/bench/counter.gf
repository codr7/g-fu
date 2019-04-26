(load "../lib/all.gf")

(class Counter ()
  ((n 0))
  (inc () (inc n))
  (dec () (dec n)))

(dump, bench 10 (for 1000
  (let c (Counter 'new))
  (c 'inc)
  (c 'dec)))

