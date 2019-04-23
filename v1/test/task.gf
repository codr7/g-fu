(let (c (chan 1))
  (push c 42)
  (test (= (len c) 1))
  (test (= (pop c) 42)))

(let (t (task () 'foo))
  (test (= (wait t) 'foo)))

(let (t1 (task () 35)
      t2 (task () 7))
  (test (= (+ (wait t1 t2)..) 42)))

(let (v 42
      t (task () (inc v)))
  (test (= (wait t) 43))
  (test (= v 42)))

(let (t (task (0 F)
          (post (fetch) 'foo)
          'bar))
  (post t (this-task))
  (test (= (fetch) 'foo))
  (test (= (wait t) 'bar)))
