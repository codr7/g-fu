(let _
  (pun foo () 35)
  (pun bar () (+ (foo) 7))
  (test (= (bar) 42)))

(let _
  (pun foo() (say "Not allowed"))

  (catch ((_ (restart 'ignore)))
    (try ((ignore () 42)) (foo))))