(debug)

(load "../lib/all.gf")

(test (= (__ foo bar baz) _))

(test (T? T))
(test (not (T? F)))

(test (= (type 42) Int))

(test (= (- (+ 1 2) (+ 3 4)) -4))

(test (= (inc 41) 42))
(test (= (inc 43 -1) 42))

(test (= (let foo 1 baz 2) 2))

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

(test (= (vec (pop-key '(foo 1 bar 2 baz 3) 'bar)) '(2 (foo 1 baz 3))))

(test (= (push () 1 2 3) '(1 2 3)))
(test (= (push '(1 2) 3) '(1 2 3)))

(let (v ())
  (push v 'foo 'bar 'baz)
  (test (= (len v) 3))
  (test (= (pop v) 'baz)))

(let (v ())
  (push v 1)
  (push v 2 3)
  (test (= (len v) 3))
  (test (= (pop v) 3))
  (test (= (peek v) 2)))

(test (= (do 1 2 3) 3))

(test (= (call (@ (fun (x) (+ x 1)) (fun (x) (* x 2))) 20) 42))

(let (fs (vec (fun (x) (+ x 1))
              (fun (x) (* x 2))))
  (test (= (call (@@ fs..) 20) 42)))

(test (= (call (fun () 42)) 42))
(test (= (call (fun (xs..) xs) 1 2 3) (vec 1 2 3)))
(test (= (call (fun (xs..) (+ xs..)) 1 2 3) 6))
(test (= (let (x 35) (call (fun (y) (+ x y)) 7)) 42))

(let _
  (fun foo ((x 42)) x)
  (test (= (foo) 42))
  (test (= (foo 7) 7)))

(let _
  (mac foo () ''bar)
  (test (= (foo) 'bar)))

(let (foo 42)
  (mac bar () 'foo)
  (test (= (bar) 42)))

(let _
  (mac foo (x) '(+ %x 7))
  (test (= (foo 35) 42)))

(let (foo 42)
  (test (= (eval 'foo) 42)))

(let (foo 35)
  (test (= (eval '(+ %foo 7)) 42)))

(let (foo (vec 35 7))
  (test (= (eval '(+ %foo..)) 42)))

(test (= (expand -1 '(foo 42)) '(foo 42)))

(let _
  (mac foo (x) x)
  (test (= (expand -1 '(foo 42)) 42)))

(let _
  (mac foo (x) x)
  (mac bar (x) '(foo %x))
  (test (= (expand 1 '(bar 42)) '(foo 42)))
  (test (= (expand 2 '(bar 42)) 42)))

(load "type.gf")
(load "cond.gf")
(load "env.gf")
(load "math.gf")
(load "iter.gf")
(load "bin.gf")
(load "fos.gf")
(load "task.gf")