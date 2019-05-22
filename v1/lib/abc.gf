(mac env (id vars body..)
  '(let %id (let %vars %body.. this-env)))

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
    (let v (new-sym) h (head cs) tcs (tail cs))
    '(let (%v %h) (if %v %(if tcs (rec tcs) v) %v)))
    
  (rec conds))

(mac or (conds..)
  (fun rec (cs)
    (let v (new-sym) h (head cs) tcs (tail cs))
    '(let (%v %h) (if %v %v %(if tcs (rec tcs) v))))
    
  (rec conds))

(mac dec (val (d 1))
  '(inc %val %(- d)))

(fun min (vals..)
  (tr (tail vals) (head vals) (fun (acc v) (if (< v acc) v acc))))

(fun max (vals..)
  (tr (tail vals) (head vals) (fun (acc v) (if (> v acc) v acc))))

(fun splat (args) args..)