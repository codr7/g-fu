(test (= (tr '(1 2 3) 0 +) 6))

(let (v '(1 2 3 4 5)
      t1 (tfilt (fun (x) (< x 4)))
      t2 (tmap (fun (x) (+ x 42)))
      ts (t@ push t1 t2)
      rts (t@@ push (vec t1 t2)..))
  (test
    (= (tr v 0 (t1 +)) 6)
    (= (tr v () (t2 push)) '(43 44 45 46 47))
    (= (tr v () ts) '(43 44 45))
    (= (tr v () rts) '(43 44 45))))

(let (v '(1 2 3 4 5)
      t (tmap (fun (x) (+ x 42)) push))
  (test (= (tr v () t) '(43 44 45 46 47))))

(let (v '(1 2 ((3 4) 5)))
  (test (= (tr v () (tcat push)) '(1 2 (3 4) 5)))
  (test (= (tr v () (tflat push)) '(1 2 3 4 5))))

(test (= (t@ 41 (fun (x) (inc x))) 42))

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

(let (v '(1 2 3))
  (test (= (find-if v (fun (x) (if (= x 2) 'ok))) 'ok))
  (test (_? (find-if v (fun (x) F)))))

(let (v '(1 2 3))
  (test (= (v/map inc) '(2 3 4))))

(let _
(fun bubbles (vs)
  (let done? F n (len vs))
  
  (while (not done?)
    (set done? T n (- n 1))
    
    (for (n i)
      (let x (# vs i) j (+ i 1) y (# vs j))
      (if (> x y) (set done? F (# vs i) y (# vs j) x))))

  vs)

(test (= (bubbles '(3 1 2)) '(1 2 3))))