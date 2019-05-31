(test (= (try _ 42) 42))

(test (= (try ((foo (bar) (+ bar 7)))
           (restart 'foo 35)
           'baz)
         42))

(test (= (try ((foo (x) x))
           (try _
             (restart 'foo 42)
             'bar)
           'baz)
         42))

(test (= (try ((foo () 'bar))
           (try ((foo () 'baz))
             (restart 'foo)))
         'baz))

(test (= (catch ((_ (restart 'done 21)))
           (try ((done (n) (* 2 n)))
             (throw 'foo)))
         42))

(test (= (catch (((Int e) (restart 'done e)))
           (try ((done (n) (* 2 n)))
             (throw 21)))
         42))
