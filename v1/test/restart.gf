(test (= (try _ 42) 42))

(test (= (try ((foo (x) (+ x 35)))
           (call (restart 'foo) 7))
         42))

(test (= (try ((foo (x) (+ x 35)))
           (try _
             (call (restart 'foo) 7)))
         42))

(test (= (try ((foo () 'bar))
           (try ((foo () 'baz))
             (call (restart 'foo))))
         'baz)) 

(test (= (catch ((_ (restart 'done 21)))
           (try ((done (n) (* 2 n)))
             (throw 'foo)))
         42))

(test (= (catch (((Int e) (restart 'done e)))
           (try ((done (n) (* 2 n)))
             (throw 21)))
         42))
