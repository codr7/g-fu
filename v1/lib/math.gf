(let z? (fun (n) (= n 0)))

(let even? (fun (n) (z? (mod n 2))))
(let odd? (fun (n) (not (even? n))))
