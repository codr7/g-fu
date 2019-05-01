(mac loop (body..)
  (let done? (new-sym) result (new-sym))
  
  '(let _
     (mac break (args..) '(recall T %args..))
     (mac continue () '(recall))
     
     ((fun ((%done? F) %result..)
        (if %done? %result.. (do %body.. (recall)))))))

(mac while (cond body..)
  '(loop
     (if %cond _ (break))
     %body..))

(mac for (args body..)
  (let v? (= (type args) Vec)
       in (new-sym)
       out (if (and v? (> (len args) 1)) (pop args) (new-sym)))
       
  '(let (%in (iter %(if v? (pop args) args)))
     (loop
       (let %out (pop %in))
       (if (_? %out) (break))
       %body..)))

(fun @ (rf fs..)
  (if (_? rf)
    (fun (rf)
      (tr (reverse fs) rf (fun (acc x) (x acc))))
    (tr (reverse fs) rf (fun (acc x) (x acc)))))

(fun tmap (f (rf _))
  (if rf
    (fun (acc val)
      (rf acc (f val)))
    (fun (rf)
      (fun (acc val)
        (rf acc (f val))))))

(fun tcat ((rf _))
  (if rf
    (fun (acc val)
      (rf acc val..))
    (fun (rf)
      (fun (acc val)
        (rf acc val..)))))

(fun tfilt (f (rf _))
  (if rf
    (fun (acc val)
      (if (f val)
        (rf acc val)
        acc))
    (fun (rf)
      (fun (acc val)
        (if (f val)
          (rf acc val)
          acc)))))