(let loop (mac (body..)
  (let done? (new-sym) result (new-sym))
  
  '(let (break (mac (args..) '(recall T %args..)))
     ((fun ((%done? F) %result..)
        (if %done? %result.. (do %body.. (recall))))))))

(let while (mac (cond body..)
  '(loop
     (if %cond _ (break))
     %body..)))

(let _? (fun (x) (= x _)))

(let for (mac (args body..)
  (let v? (= (type args) Vec)
       in (new-sym)
       out (if (and v? (> (len args) 1)) (pop args) (new-sym)))
       
  '(let (%in (iter %(if v? (pop args) args)))
     (loop
       (let %out (pop %in))
       (if (_? %out) (break))
       %body..))))

(let @ (fun (rf fs..)
  (if (= rf _)
    (fun (rf)
      (fold (reverse fs) rf (fun (acc x) (x acc))))
    (fold (reverse fs) rf (fun (acc x) (x acc))))))

(let map (fun (f (rf _))
  (if rf
    (fun (acc val)
      (rf acc (f val)))
    (fun (rf)
      (fun (acc val)
        (rf acc (f val)))))))

(let cat (fun (rf)
  (fun (acc val)
    (push acc val..))))

(let keep (fun (f (rf _))
  (if rf
    (fun (acc val)
      (if (f val)
        (rf acc val)
        acc))
    (fun (rf)
      (fun (acc val)
        (if (f val)
          (rf acc val)
          acc))))))