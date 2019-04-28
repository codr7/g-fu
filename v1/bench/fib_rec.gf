(load "../lib/core.gf")
(load "../lib/iter.gf")

(let (fib (fun n
            (if (< n 2)
              n
              (+ (fib (- n 1)) (fib (- n 2))))))
  (dump (bench 10 (for 10 (fib 20)))))