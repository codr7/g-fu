(test (= (Int/? 42) Int))
(test (_? (Int/? T)))
(test (= (Seq/? Vec) Seq))
(test (= (Seq/? IntIter) Iter))