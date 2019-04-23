(load "../lib/core.gf")
(load "../lib/iter.gf")

(let (fib (fun (n a b)
            (if n 
              (if (= n 1)
                b
                (recall (- n 1) b (+ a b)))
              a)))
  (dump (bench 10 (for 10000 (fib 20 0 1)))))