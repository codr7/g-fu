(test (= (try _ 42) 42))

(try ((foo (bar) (+ bar 7)))
  (test (= (restart 'foo 35) 42)))

