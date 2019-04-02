(test (as-bool 42))
(test (not (as-bool 0)))

(test (= (and T 42) 42))
(test (= (or 42 T _) 42))

(test (= (- 42) -42))
(test (= (- 10 1 2 3) 4))

(test (= (+ -42) 42))
(test (= (+ 35 7) 42))

(test (= (Vec 1 2 3) (Vec 1 2 3)))
(test (not (== (Vec 1 2 3) (Vec 1 2 3))))
(test (= (+ (Vec 1 2 3)..) 6))

(test (= (do 1 2 3) 3))

(test (= ((fun (xs..) xs) 1 2 3) (Vec 1 2 3)))
(test (= ((fun (xs..) (+ xs..)) 1 2 3) 6))
(test (= (let (x 35) ((fun (y) (+ x y)) 7)) 42))