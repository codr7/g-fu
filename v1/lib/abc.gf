(let NOP (mac (args..)))

(fun tr (in acc fn)
  (let rec (fun (in acc fn)
    (let v (pop in))
    (if (_? v) acc (recall in (fn acc v) fn))))

  (rec (iter in) acc fn))

(fun ~ (fs..)
  (tr (reverse fs) _ (fun (acc x) (x acc))))

(fun not (val)
  (if val F T))

(let and (mac (conds..)
  (fun rec (cs)
    (let h (head cs) tcs (tail cs))
    '(if %h %(if tcs (rec tcs) h)))
    
  (rec conds)))

(let or (mac (conds..)
  (fun rec (cs)
    (let h (head cs) tcs (tail cs))
    '(if %h %h %(if tcs (rec tcs))))
    
  (rec conds)))
  
(let dec (mac (var (delta 1))
  '(inc %var (- %delta))))

(fun splat (args) args..)