(debug)

(load "../lib/abc.gf")
(load "../lib/iter.gf")

(fun Widget (args..)
  (use _ do)
  (let left 0 top 0
       width (or (pop-key args 'width) (fail "Missing width"))
       height (or (pop-key args 'height) (fail "Missing height")))
  (this-env))

(fun Button (args..)
  (let Widget (Widget args..))
  (Widget/do
    (this-env)))

(let (w (Widget 'width 100 'height 50))
  (dump w/width))

(let (b (Button 'width 100 'height 50))
  (dump b/width))

(__ (b (Button/new 'width 100 'height 50))
  (test (= (b/move 10 10) '(10 10)))
  (test (= (b/resize 0 0) '(142 50)))

  (let (called F)
    (b/on-click (fun (b) (set 'called T)))
    (test (not called))
    (b/click)
    (test called)))
