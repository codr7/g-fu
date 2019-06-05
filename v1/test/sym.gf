(test (== 'foo 'foo))
(test (= ''foo ''foo))

(test (= $foo $foo))
(test (not (= (let _ $foo) (let _ $foo))))