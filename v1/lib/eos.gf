(debug)

(load "../lib/abc.gf")
(load "../lib/iter.gf")

(let Widget (let _
  (fun new (args..)
    (let left 0 top 0
         width (or (pop-key args 'width) (fail "Missing width"))
         height (or (pop-key args 'height) (fail "Missing height")))

    (fun move (dx dy)
      (vec (inc left dx)
           (inc top dy)))

    (fun resize (dx dy)
      (vec (inc width dx)
           (inc height dy)))
  
    this-env)

  this-env))

(let Button (let _
  (fun new (args..)
    (let w (Widget/new args..))
    (use w move)
  
    (fun resize (dx dy)
      (w/resize (+ dx 42) dy))
    
    this-env)

  this-env))

(dump Button)

(let (b (Button/new 'width 100 'height 50))
  (dump b)
  (test (= (b/resize 0 0) '(142 50)))
  (test (= (b/move 10 10) '(10 10)))

  (__ (called F)
    (b/on-click (fun (b) (set 'called T)))
    (test (not called))
    (b/click)
    (test called)))