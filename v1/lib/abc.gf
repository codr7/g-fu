(mac NOP (args..))

(fun tr (in acc fn)
  (let rec (fun (in acc fn)
    (let v (pop in))
    (if (_? v) acc (recall in (fn acc v) fn))))

  (rec (iter in) acc fn))

(fun ~ (fs..)
  (tr (reverse fs) _ (fun (acc x) (x acc))))

(fun not (val)
  (if val F T))

(mac and (conds..)
  (fun rec (cs)
    (let h (head cs) tcs (tail cs))
    '(if %h %(if tcs (rec tcs) h)))
    
  (rec conds))

(mac or (conds..)
  (fun rec (cs)
    (let h (head cs) tcs (tail cs))
    '(if %h %h %(if tcs (rec tcs))))
    
  (rec conds))
  
(mac dec (var (delta 1))
  '(inc %var (- %delta)))

(fun splat (args) args..)