(let (c (chan 1))
  (push c 42)
  
  (test (= (len c) 1))
  (test (= (pop c) 42)))

(test (= (wait (task () 'foo)) 'foo))

(let _
  (task t1 () 35)
  (task t2 () 7)
  
  (test (= (+ (wait t1 t2)..) 42)))

(let (v 42)
  (task t () (inc v))  
  (test (= (wait t) 43))
  (test (= v 42)))

(let _
  (task t (0)
    (post (fetch) 'foo)
    'bar)
    
  (post t this-task)
  (test (= (fetch) 'foo))
  (test (= (wait t) 'bar)))
