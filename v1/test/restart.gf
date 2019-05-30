(test (= (try _ 42) 42))

(test (= (try ((foo (bar) (+ bar 7)))
           (restart 'foo 35)
           'baz)
         42))

(test (= (try ((foo (bar) bar))
           (try _
             (restart 'foo 42)
             'baz)
           'qux)
         42))
