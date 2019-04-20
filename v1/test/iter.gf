(load "../lib/iter.gf")

(test (= (loop (break 'foo)) 'foo))

(let (i 0)
  (while (< (inc i) 7))
  (test (= i 7)))

(let (i 0)
  (for 7 (inc i))
  (test (= i 7)))

(let (i 0)
  (for (7 j) (inc i j))
  (test (= i 21)))

(let (v '(1 2 3 4 5)
      t1 (@ (keep (fun (x) (< x 4))))
      t2 (@ (map (fun (x) (+ x 42)))))
  (test (= (fold v (t1 +) 0) 6))
  (test (= (fold v (t2 push)) '(43 44 45 46 47)))
  (test (= (fold v (t1 (t2 push)) '(43 44 45)))))