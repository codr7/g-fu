(test (= (len (let _ (this-env))) 1))

(let (e (let (foo 1)
          (use / do let)
          (this-env)))
  (test (= e/foo 1))
  (e/do (let bar 2))
  (test (= e/bar 2)))