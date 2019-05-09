(test (= (len (let _ this-env)) 1))

(let (e (let (foo 1)
          (use _ do let)
          this-env))
  (test (= e/foo 1))
  (e/do (let bar 2))
  (test (= e/bar 2))

  (let ee (let _ (use e foo))
    (test (= ee/foo 1))))

(let (foo (let (bar 7)
            (fun resolve (id) 42)
            this-env))
  (test (= foo/bar 7))
  (test (= foo/baz 42)))