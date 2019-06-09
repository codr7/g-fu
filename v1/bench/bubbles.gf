(load "../lib/all.gf")

(fun bubbles (vs)
  (let done? F n (len vs))
  
  (while (not done?)
    (set done? T n (- n 1))
    
    (for (n i)
      (let x (# vs i) j (+ i 1) y (# vs j))
      (if (> x y) (set done? F (# vs i) y (# vs j) x))))

  vs)

(let vals ())

(for 100
  (push vals (rand 100000)))

(dump (bench 10
  (bubbles (clone vals))))

