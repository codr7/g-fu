(load "../lib/fos.gf")

(let (n 0 fo (dispatch
              (inc ((d 1)) (inc n d))
              (dec ((d 1)) (dec n d))))
  (test (= (fo 'inc 4) 4))
  (test (= (fo 'dec) 3))
  (test (= n 3)))

(let-self ()
  (test (= self 42))
  42)

(let (fo (let-self ()
           (dispatch
             (patch (new) (set self new)))))
  (fo 'patch (fun (x) x))
  (test (= (fo 42) 42)))