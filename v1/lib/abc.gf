(mac NOP (args..))

(fun tr (in acc fn)
  (fun rec (in acc fn)
    (let v (pop in))
    (if (_? v) acc (rec in (fn acc v) fn)))

  (rec (iter in) acc fn))

(mac @ (f1 fs..)
  '(fun (args..)
     %(tr fs '(call %f1 args..) (fun (acc x) '(call %x %acc)))))

(fun @@ (f1 fs..)
  (fun (args..)
    (tr fs (f1 args..) (fun (acc x) (x acc)))))
  
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