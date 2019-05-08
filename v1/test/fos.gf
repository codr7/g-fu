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
             (patch (new) (set 'self new)))))
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
  (on-click)

  (resize (dx dy)
    (self 'Widget/resize (+ dx 42) dy))

  (on-click (f)
    (push on-click f))

  (click ()
    (for (on-click f) (f self))))

(let b (Button 'new 'width 100 'height 50))
(test (= (b 'move 10 10) '(10 10)))
(test (= (b 'resize 0 0) '(142 50)))

(let (called F)
  (b 'on-click (fun (b) (set 'called T)))
  (test (not called))
  (b 'click)
  (test called))
  
)
