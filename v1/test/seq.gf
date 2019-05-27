(let (v '(foo bar baz))
  (test (= (v/join _) "foobarbaz"))
  (test (= (v/join ",") "foo,bar,baz")))