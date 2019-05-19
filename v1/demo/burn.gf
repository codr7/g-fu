(load "../lib/all.gf")

(let rows 24 cols 25
     buf (new-bin (* rows cols)))

(fun offs (r c)
  (+ (* (- rows r 1) cols) (- cols c 1)))

(fun render ()
  (for (rows r)
    (for (cols c)
      (let i (offs r c))
      (set (# buf i) 0xff))))
      
(render)