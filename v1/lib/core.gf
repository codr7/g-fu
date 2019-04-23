(let @ (fun (fs..)
  (fun (in)
    (fold (reverse fs) (fun (acc x) (x acc)) in))))

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