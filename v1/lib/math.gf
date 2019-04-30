(let z? (fun (n) (= n 0)))
(let even? (fun (n) (z? (mod n 2))))
(let odd? (fun (n) (not (even? n))))

(let exp (fun (base n)
  (switch
    ((z? n) 1)
    ((even? n) (* (exp base (div n 2))))
    (T (* base (exp base (- n 1)))))))

(let fib (fun (n (a 0) (b 1))
  (if n 
    (if (= n 1)
      b
      (recall (- n 1) b (+ a b)))
    a)))