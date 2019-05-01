(let NOP (mac (args..)))

(let tr (fun (in acc fn)
  (let rec (fun (in acc fn)
    (let v (pop in))
    (if (_? v) acc (recall in (fn acc v) fn))))

  (rec (iter in) acc fn)))

(let ~ (fun (fs..)
  (tr (reverse fs) _ (fun (acc x) (x acc)))))

(let not (fun (val)
  (if val F T)))

(let and (mac (conds..)
  (let rec (fun (cs)
    (let h (head cs) tcs (tail cs))
    '(if %h %(if tcs (rec tcs) h))))
  (rec conds)))

(let or (mac (conds..)
  (let rec (fun (cs)
    (let h (head cs) tcs (tail cs))
    '(if %h %h %(if tcs (rec tcs)))))
  (rec conds)))
  
(let dec (mac (var (delta 1))
  '(inc %var (- %delta))))

(let splat (fun (args) args..))