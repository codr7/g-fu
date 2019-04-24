(load "../lib/fos.gf")

(let (n 0 d (dispatch
              (inc ((delta 1)) (inc n delta))
              (dec ((delta 1)) (dec n delta))))
  (test (= (d 'inc 4) 4))
  (test (= (d 'dec) 3))
  (test (= n 3)))

(let-self ()
  (test (= self 42))
  42)

(let (s (let-self ()
           (dispatch
             (patch (new) (set self new)))))
  (s 'patch (fun (x) x))
  (test (= (s 42) 42)))