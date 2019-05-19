(let rows 25
     cols 25
     buf (new-bin (* rows cols)))

(fun xy (x y)
  (# buf (+ (* y cols) x)))

(fun set-xy (f x y)
  (let i (+ (* y cols) x))
  (set (# buf i) (f (# buf i)))
  _)

(set (xy 10 10) 0xff)
(dump (xy 10 10))