(load "../lib/fos.gf")

(let (n 0 c (funs
              (inc ((d 1)) (inc n d))
              (dec ((d 1)) (dec n d))))
  (test (= (c 'inc 4) 4))
  (test (= (c 'dec) 3))
  (test (= n 3)))