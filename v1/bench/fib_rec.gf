(load "../lib/all.gf")

(fun f (n)
  (if (< n 2)
    n
    (+ (f (- n 1)) (f (- n 2)))))
    
(dump (bench 10 (for 10 (f 20))))