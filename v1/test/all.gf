(test (Bool 42))
(test (not (Vec)))

(test (= (and T 42) 42))
(test (= (or 42 T _) 42))

(let (x 35)
  (test (= (inc x 7) 42))
  (test (= x 42)))

(test (= (- 42) -42))
(test (= (- 10 1 2 3) 4))

(test (= (+ -42) 42))
(test (= (+ 35 7) 42))

(test (== 'foo 'foo))
(test (= ''foo ''foo))
(test (not (= (Sym) (Sym))))

(test (= '(1 2 3) (Vec 1 2 3)))
(test (= '(1 %(+ 2 3) 4) (Vec 1 5 4)))
(test (= (Vec 1 2 3) (Vec 1 2 3)))
(test (not (== (Vec 1 2 3) (Vec 1 2 3))))
(test (= '(1 %(Vec 2 3)..) (Vec 1 2 3)))
(test (= (+ (Vec 1 2 3)..) 6))

(let (v (Vec))
  (push v 1)
  (push v 2 3)
  (test (= (len v) 3))
  (test (= (pop v) 3))
  (test (= (peek v) 2)))

(test (= (do 1 2 3) 3))

(test (= ((fun (xs..) xs) 1 2 3) (Vec 1 2 3)))
(test (= ((fun (xs..) (+ xs..)) 1 2 3) 6))
(test (= (let (x 35) ((fun (y) (+ x y)) 7)) 42))

(let (foo (fun (x?) x))
  (test (= (foo 42) 42))
  (test (= (foo) _)))

(let (foo (macro () ''bar))
  (test (= (foo) 'bar)))

(let (foo 42 bar (macro () 'foo))
  (test (= (bar) 42)))

(let (foo (macro (x) '(+ %x 7)))
  (test (= (foo 35) 42)))

(let (fib (fun (n)
            (if (< n 2)
              n
              (+ (fib (- n 1)) (fib (- n 2))))))
  (test (= (fib 20) 6765)))

(let (fib (fun (n a b)
            (if n 
              (if (= n 1)
                b
                (recall (- n 1) b (+ a b)))
              a)))
  (test (= (fib 20 0 1) 6765)))


(let (foo 42)
  (test (= (eval 'foo) 42)))

(let (foo 35)
  (test (= (eval '(+ %foo 7)) 42)))

(let (foo (Vec 35 7))
  (test (= (eval '(+ %foo..)) 42)))