(mac env (id vars body..)
  '(let %id (let %vars %body.. Env/this)))

(fun tr (in acc fn)
  (fun rec (in acc)
    (let v (pop in))
    (if (_? v) acc (recall in (fn acc v))))

  (rec (in/iter) acc))

(mac @ (f1 fs..)
  '(fun (args..)
     %(tr fs '(call %f1 args..) (fun (acc x) '(call %x %acc)))))

(fun @@ (f1 fs..)
  (fun (args..)
    (tr fs (f1 args..) (fun (acc x) (x acc)))))
  
(pun not (val)
  (if val F T))

(mac and (conds..)
  (pun rec (cs)
    (let h (head cs) tcs (tail cs))
    '(let (%$v %h) (if %$v %(if tcs (rec tcs) $v) %$v)))
    
  (rec conds))

(mac or (conds..)
  (pun rec (cs)
    (let h (head cs) tcs (tail cs))
    '(let (%$v %h) (if %$v %$v %(if tcs (rec tcs) $v))))
    
  (rec conds))

(mac dec (val (d 1))
  '(inc %val %(- d)))

(fun min (vals..)
  (tr (tail vals) (head vals) (pun (acc v) (if (< v acc) v acc))))

(fun max (vals..)
  (tr (tail vals) (head vals) (pun (acc v) (if (> v acc) v acc))))

(pun splat (args) args..)