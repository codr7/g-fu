(fun tr (in acc fn)
  (fun rec (in acc fn)
    (let v (pop in))
    (if (_? v) acc (rec in (fn acc v) fn)))

  (rec (iter in) acc fn))

(mac use (prefix ids..)
  (if (_? prefix)
    '(__ (use _ %ids..))
    '(let %(tr ids _
               (fun (acc s)
                 (push acc s (sym prefix '/ s))))..)))

(mac @ (f1 fs..)
  '(fun (args..)
     %(tr fs '(call %f1 args..) (fun (acc x) '(call %x %acc)))))

(fun @@ (f1 fs..)
  (fun (args..)
    (tr fs (f1 args..) (fun (acc x) (x acc)))))
  
(fun not (val)
  (if val F T))

(mac and (conds..)
  (let v (new-sym))
  
  (fun rec (cs)
    (let h (head cs) tcs (tail cs))
    '(let (%v %h) (if %v %(if tcs (rec tcs) v))))
    
  (rec conds))

(mac or (conds..)
  (let v (new-sym))
  
  (fun rec (cs)
    (let h (head cs) tcs (tail cs))
    '(let (%v %h) (if %v %v %(if tcs (rec tcs)))))
    
  (rec conds))
  
(mac dec (var (delta 1))
  '(inc %var (- %delta)))

(fun splat (args) args..)