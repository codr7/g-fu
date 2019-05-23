(debug)
(load "../lib/all.gf")

(env fire (width 50 height 25
           buf (new-bin (* width height))
           out stdout
           max-fade 50
           tot-frames 0 tot-time .0)
  (fun ctrl (args..)
    (print out "\e[" args..))

  (fun clear ()
    (ctrl "2J"))

  (fun home ()
    (ctrl "H"))

  (fun init ()
    (for (width i)
      (set (# buf i) 0xff))

    (clear))

  (fun render ()
    (let t0 (now) i -1)

    (for ((- height 1) y)
      (for (width x)
        (let v (# buf (inc i))
             j (+ i width))
        
        (if (and x (< x (- width 1)))
          (inc j (- 1 (rand 3))))
        
        (set (# buf j)
             (if v (- v (rand (min max-fade (+ (int v) 1)))) v))))

    (inc i (+ width 1))
    (let prev-g _)
    (home)
    
    (for (height y)
      (for (width x)
        (let g (# buf (dec i))
             r (if g 0xff 0)
             b (if (= g 0xff) 0xff 0))
             
        (if (= g prev-g)
          (print out " ")
          (do
            (ctrl "48;2;" (int r) ";" (int g) ";" (int b) "m ")
            (set prev-g g))))

      (print out \n))

    (flush out)
    (inc tot-time (- (now) t0))
    (inc tot-frames))

  (fun restore ()
    (ctrl "0m")
    (clear)
    (home)))

(fire/init)
(for 50 (fire/render))
(fire/restore)

(say (/ (* 1000000000.0 fire/tot-frames) fire/tot-time))