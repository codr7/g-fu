(load "../doc/fos.gf")

(let (n 0 d (dispatch
              (inc ((delta 1)) (inc n delta))))
  (test (= (d 'inc 4) 4))
  (test (= (d 'inc -1) 3))
  (test (= n 3)))

(let-this ()
  (test (= this 42))
  42)

(let (s (let-this ()
           (dispatch
             (patch (new) (set this new)))))
  (s 'patch (fun (x) x))
  (test (= (s 42) 42)))

(let _

(class Widget ()
  ((left 0) (top 0)
   (width (fail "Missing width")) (height (fail "Missing height")))
  
  (move (dx dy)
    (vec (inc left dx)
         (inc top dy)))

  (resize (dx dy)
    (vec (inc width dx)
         (inc height dy))))

(class Button (Widget)
  ((on-click ()))

  (resize (dx dy)
    (this 'Widget/resize (+ dx 42) dy))

  (on-click (f)
    (push on-click f))

  (click ()
    (for (on-click f) (f this))))

(let b (Button 'new 'width 100 'height 50))
(test (= (b 'move 10 10) '(10 10)))
(test (= (b 'resize 0 0) '(142 50)))

(let (called F)
  (b 'on-click (fun (b) (set called T)))
  (test (not called))
  (b 'click)
  (test called))
  
)
