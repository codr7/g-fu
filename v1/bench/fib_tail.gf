(load "../lib/all.gf")

(let f (fun (n (a 0) (b 1))
  (if n 
    (if (= n 1)
      b
      (recall (- n 1) b (+ a b)))
    a)))
    
(dump (bench 10 (for 10000 (f 20))))