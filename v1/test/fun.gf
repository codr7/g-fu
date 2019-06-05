(test (= (type (fun ())) Fun))
(test (= (type (pun ())) Pun))

(let _
  (pun foo() 35)
  (pun bar() (+ (foo) 7))
  (test (= (bar) 42)))

(let _
  (pun foo() (say "Not allowed"))

  (test (= (catch (((EImpure _) (restart 'ignore)))
             (try ((ignore () 'ok)) (foo)))
           'ok)))

(let _
  (env foo (x ()))
  (pun bar() foo/x)
  (push (bar) 42)
  (test (= foo/x ())))

(let _
  (env foo (x 42))
  (pun bar() (inc foo/x))

  (test (= (catch (((EImpure _) (restart 'ignore)))
             (try ((ignore () 'ok)) (bar)))
           'ok)))
