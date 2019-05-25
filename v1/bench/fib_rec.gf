(load "../lib/all.gf")

(fun f (n)
  (if (< n 2)
    n
    (+ (f (dec n)) (f (dec n)))))
    
(dump (bench 10 (for 10 (f 20))))