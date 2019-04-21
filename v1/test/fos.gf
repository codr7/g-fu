(load "../lib/fos.gf")

(let (n 0 fo (fo-fun
              (inc ((d 1)) (inc n d))
              (dec ((d 1)) (dec n d))))
  (test (= (fo 'inc 4) 4))
  (test (= (fo 'dec) 3))
  (test (= n 3)))

(fo-let ()
  (test (= self 42))
  42)