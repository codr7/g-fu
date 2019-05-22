(test (z? 0))
(test (not (z? 42)))
(test (not (z? -42)))

(test (even? 42))
(test (not (even? 21)))

(test (odd? 21))
(test (not (odd? 42)))

(test (= (gcd 14 21) 7))

(let _
  (fun f (n)
    (if (< n 2)
      n
      (+ (f (- n 1)) (f (- n 2)))))
  
  (test (= (f 20) 6765)))

(test (= (fib 20 0 1) 6765))

(test (= (+ .5 .25) .75))
(test (= (+ -.5) .5))
(test (= (+ .33 .33) .66))
(test (= (+ (/ 1 3) (/ 2 3)) 1.))

(test (Float 42.))

(test (= (+ (float 42) .5) 42.5))