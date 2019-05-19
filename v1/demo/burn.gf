(load "../lib/all.gf")

(let width 25 height 24
     buf (new-bin (* width (+ height 1))))

(fun offs (x y)
  (+ (- width x 1) (* (- height y 1) width)))

(fun xy (x y)
  (# buf (offs x y)))

(fun set-xy (f x y)
  (let i (offs x y))
  (set (# buf i) (f (# buf i))))

(fun render ()
  (for (height y)
    (for (width x)
      (set (xy x y) 0xff))))
      
(render)