(load "../lib/iter.gf")

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
      t1 (keep (fun (x) (< x 4)))
      t2 (map (fun (x) (+ x 42)))
      ts (@ push t1 t2))
  (test
    (= (fold v 0 (t1 +)) 6)
    (= (fold v _ (t2 push)) '(43 44 45 46 47))
    (= (fold v _ ts) '(43 44 45))))

(let (v '(1 2 3 4 5)
      t (map (fun (x) (+ x 42)) push))
  (test (= (fold v _ t) '(43 44 45 46 47))))

(let (v '((1 2) (3 4) (5)))
  (test (= (fold v _ (cat push)) '(1 2 3 4 5))))

(test (= (@ 41 (fun (x) (inc x))) 42))
(test (= ((@ _ (fun (x) (inc x))) 41) 42))