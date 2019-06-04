(let _
  (pun foo() 35)
  (pun bar() (+ (foo) 7))
  (test (= (bar) 42)))

(let _
  (pun foo() (say "Not allowed"))

  (test (= (catch ((_ (restart 'ignore)))
             (try ((ignore () 42)) (foo)))
           42)))

(let _
  (pun foo() (fun bar()))

  (test (= (catch (((EImpure _) (restart 'ignore)))
             (try ((ignore() 42))
               (foo)
               (fail "Not reached"))))))