(load "../lib/core.gf")
(load "../lib/iter.gf")

(dump (bench 10
  (let (s _)
    (for (100000 i) (push s i))
    (for 100000 (pop s)))))