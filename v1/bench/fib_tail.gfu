(let (fib (fun (n a b)
            (if (z? n)
              a
              (if (one? n)
                b
                (recall (- n 1) b (+ a b))))))
  (dump (bench 10 (for 10000 (fib 20 0 1)))))