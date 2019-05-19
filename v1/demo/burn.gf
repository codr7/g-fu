(load "../lib/all.gf")

(let rows 25 cols 25
     buf (new-bin (* rows cols)))

(fun get-offs (x y)
  (+ (* y cols) x))

(fun render ()
  (for (rows r)
    (for (cols c)
      (let i (get-offs r c))
      (set (# buf i) 0xff))))
      
(render)