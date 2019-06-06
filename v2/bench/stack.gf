(load "../lib/all.gf")

(dump (bench 10
  (let (s ())
    (for (100000 i) (push s i))
    (for 100000 (pop s)))))