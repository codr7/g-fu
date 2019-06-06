(pun z? (n) (= n 0))

(pun even? (n) (z? (mod n 2)))

(pun odd? (n) (not (even? n)))

(pun exp (base n)
  (switch
    ((z? n) 1)
    ((even? n) (* (exp base (div n 2))))
    (T (* base (exp base (- n 1))))))

(pun gcd (a b)
  (if b
    (recall b (mod a b))
    a))

(pun fib (n (a 0) (b 1))
  (if n 
    (if (= n 1)
      b
      (recall (- n 1) b (+ a b)))
    a))