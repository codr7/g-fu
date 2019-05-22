(debug)
(load "../lib/all.gf")

(env fire (width 50 height 25
           esc (str 0x1b "[")
           buf (new-bin (* width height 1))
           out stdout)
  (fun get-offs (x y)
    (+ (- width x 1) (* (- height y 1) width)))

  (fun xy (x y)
    (# buf (get-offs x y)))

  (fun set-xy (f x y)
    (let i (get-offs x y))
    (set (# buf i) (f (# buf i)))
    _)

  (fun move-to (x y)
    (print out (str esc (+ y 1) ";" (+ x 1) "H")))

  (fun pick-color (r g b)
    (print out (str esc "48;2;" (int r) ";" (int g) ";" (int b) "m")))

  (fun init ()
    (for (width x)
      (set (xy x 0) 0xff)))

  (fun render ()
    (for ((- height 1) y)
      (for (width x)
        (let v (xy x y))
        (set (xy x (+ y 1)) (- v (min 50 (rand (+ v 1)))))))

    (let i -1)
    
    (for (height y)
      (move-to 0 y)

      (for (width x)
        (let g (# buf (inc i)) r (if g 0xff 0x00) b (if (= g 0xff) 0xff 0x00))
        (pick-color r g b)
        (print out " ")))

    (flush out)))

(fire/init)

(for 50
  (fire/render))