(debug)
(load "../lib/all.gf")

(let width 25 height 25
     data (new-bin (* width height 1)))

(fun get-offs (x y)
  (+ (- width x 1) (* (- height y 1) width)))

(fun xy (x y)
  (# data (get-offs x y)))

(fun set-xy (f x y)
  (let i (get-offs x y))
  (set (# data i) (f (# data i)))
  _)

(let esc (str 0x1b "[")
     top-left (str esc "H"))

(fun move-top-left (out)
  (print out top-left))

(fun pick-color (out r g b)
  (print out (str esc "48;2;" r ";" g ";" b "m")))

(fun render (out)
  (for ((- height 1) y)
    (for (width x)
      (let v (xy x y))
      (set (xy x (+ y 1)) (- v (rand (+ v 1))))))

  (move-top-left out)
  
  (for (height y)
    (if y (move-next-line out))

    (for (width x)
      (let g (xy x y) r (if g 0xff 0x00) b (if (= g 0xff) 0xff 0x00))
      (pick-color out r g b)
      (print out " " ))))

(for (width x)
  (set (xy x 0) 0xff))

(render _)