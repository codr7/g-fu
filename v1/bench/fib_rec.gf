(load "../lib/all.gf")

(let f (fun n
  (if (< n 2)
    n
    (+ (f (- n 1)) (f (- n 2))))))
    
(dump (bench 10 (for 10 (f 20))))