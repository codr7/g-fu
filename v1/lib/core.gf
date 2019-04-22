(let @ (fun (fs..)
  (fun (in)
    (fold (reverse fs) (fun (acc x) (x acc)) in))))

(let dec (mac (var (delta 1))
  '(inc %var (- %delta))))

(let splat (fun (args) args..))