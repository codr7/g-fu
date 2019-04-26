(debug)

(load "../lib/core.gf")

(test (= (NOP foo bar baz) _))

(test (= (type 42) Int))

(test (= (-, + 1 2, + 3 4) -4))

(let (x 35)
  (test (= (inc x 7) 42))
  (test (= x 42)))

(test (= (- 42) -42))
(test (= (- 10 1 2 3) 4))

(test (= (+ -42) 42))
(test (= (+ 35 7) 42))

(test (== 'foo 'foo))
(test (= ''foo ''foo))
(test (not (= (new-sym) (new-sym))))

(test (= (len "") 0))
(test (= (len "foo") 3))

(test (= (len (vec)) 0))
(test (= (len '(1 2 3)) 3))

(test (= '(1 2 3) (vec 1 2 3)))
(test (= '(1 %(+ 2 3) 4) (vec 1 5 4)))
(test (= (vec 1 2 3) (vec 1 2 3)))
(test (= '(1 %(vec 2 3)..) (vec 1 2 3)))
(test (= (+ (vec 1 2 3)..) 6))

(let (v '(foo 1 bar 2 baz 3))
  (test (= (find-key v 'foo) 1))
  (test (= (find-key v 'baz) 3))
  (test (= (find-key v 'qux) _)))

(let (v '(foo 1 bar 2 baz 3))
  (test (= (pop-key v 'bar) 2))
  (test (= v '(foo 1 baz 3))))

(test (= (push _ 1 2 3) '(1 2 3)))
(test (= (push '(1 2) 3) '(1 2 3)))

(let (v _)
  (push v 'foo 'bar 'baz)
  (test (= (len v) 3))
  (test (= (pop v) 'baz)))

(let (v (vec))
  (push v 1)
  (push v 2 3)
  (test (= (len v) 3))
  (test (= (pop v) 3))
  (test (= (peek v) 2)))

(test (= (do 1 2 3) 3))

(test (= ((fun () 42)) 42))
(test (= ((fun (xs..) xs) 1 2 3) (vec 1 2 3)))
(test (= ((fun (xs..) (+ xs..)) 1 2 3) 6))
(test (= (let (x 35) ((fun (y) (+ x y)) 7)) 42))

(let (foo (fun ((x)) x))
  (test (= (foo) _))
  (test (= (foo 42) 42)))

(let (foo (fun ((x 42)) x))
  (test (= (foo) 42))
  (test (= (foo 7) 7)))

(let (foo (mac () ''bar))
  (test (= (foo) 'bar)))

(let (foo 42 bar (mac () 'foo))
  (test (= (bar) 42)))

(let (foo (mac (x) '(+ %x 7)))
  (test (= (foo 35) 42)))

(let (fib (fun n
            (if, < n 2,
              n,
              (+, fib (- n 1), fib (- n 2)))))
  (test, = (fib 20) 6765))

(NOP

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

(let (foo (vec 35 7))
  (test (= (eval '(+ %foo..)) 42)))

(test (= (expand -1 '(foo 42)) '(foo 42)))

(let (foo (mac (x) x))
  (test (= (expand -1 '(foo 42)) 42)))

(let (foo (mac (x) x)
      bar (mac (x) '(foo %x)))
  (test (= (expand 1 '(bar 42)) '(foo 42)))
  (test (= (expand 2 '(bar 42)) 42)))

(load "cond.gf")
(load "iter.gf")
(load "fos.gf")
(load "task.gf")
)