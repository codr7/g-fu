(test (z? 0))
(test (not (z? 42)))
(test (not (z? -42)))

(test (even? 42))
(test (not (even? 21)))

(test (odd? 21))
(test (not (odd? 42)))

(let (f (fun n
          (if (< n 2)
            n
            (+ (f (- n 1)) (f (- n 2))))))
  (test (= (f 20) 6765)))

(test (= (fib 20 0 1) 6765))
