(let splat (fun (args) args..))

(let @ (fun (fs..)
  (fun (in)
    (fold (reverse fs) (fun (acc x) (x acc)) in))))