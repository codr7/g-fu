(test (== 'foo 'foo))
(test (= ''foo ''foo))

(test (not (= (new-sym) (new-sym))))

(test (= $foo $foo))