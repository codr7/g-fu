(if (val 'lib-loaded?) _ (do

(let lib-loaded? T)

(load "../lib/abc.gf")
(load "../lib/iter.gf")
(load "../lib/cond.gf")
(load "../lib/math.gf")

))