(test (= (loop (break 'foo)) 'foo))

(let (i 0)
  (test
    (= (loop (inc i) (if (= i 42) (break 'done) (continue))) 'done)
    (= 42 i))) 

(let (i 0)
  (while (< (inc i) 7))
  (test (= i 7)))

(let (i 0)
  (for 7 (inc i))
  (test (= i 7)))

(let (i 0)
  (for (7 j) (inc i j))
  (test (= i 21)))

(let (v '(1 2 3) n 0)
  (for (v i) (inc n i))
  (test (= n 6)))

(let (v '(1 2 3 4 5) n 0)
  (drop v)
  (test (= v '(1 2 3 4)))
  (drop v 2)
  (test (= v '(1 2))))

(let (v '(1 2 3 4 5)
      t1 (tfilt (fun (x) (< x 4)))
      t2 (tmap (fun (x) (+ x 42)))
      ts (t@ push t1 t2))
  (test
    (= (tr v 0 (t1 +)) 6)
    (= (tr v _ (t2 push)) '(43 44 45 46 47))
    (= (tr v _ ts) '(43 44 45))))

(let (v '(1 2 3 4 5)
      t (tmap (fun (x) (+ x 42)) push))
  (test (= (tr v _ t) '(43 44 45 46 47))))

(let (v '((1 2) (3 4) (5)))
  (test (= (tr v _ (tcat push)) '(1 2 3 4 5))))

(let (out _)
  (tr 3 _ (tpipe (fun (x) (push out (+ x 1)))))
  (test (= out '(1 2 3))))

(test (= (t@ 41 (fun (x) (inc x))) 42))
