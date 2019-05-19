(load "../lib/all.gf")

(let width 25 height 25
     buf (new-bin (* width height 1)))

(fun offs (x y)
  (+ (- width x 1) (* (- height y 1) width)))

(fun xy (x y)
  (# buf (offs x y)))

(fun set-xy (f x y)
  (let i (offs x y))
  (set (# buf i) (f (# buf i)))
  _)

(fun rand (x)
  (div x 2))

(fun render ()
  (for ((- height 1) y)
    (for (width x)
      (let v (xy x y))
      (set (xy x (+ y 1)) (- v (rand v))))))

(for (width x)
  (set (xy x 0) 0xff))

(render)