(let loop (mac (body..)
  (let done? (new-sym) result (new-sym))
  
  '(let (break (mac (args..) '(recall T %args..)))
     ((fun ((%done? F) %result..)
        (if %done? %result.. (do %body.. (recall))))))))

(let while (mac (cond body..)
  '(loop
     (if %cond _ (break))
     %body..)))

(let for (mac (args body..)
  (let v? (= (type args) Vec)
       i (if (and v? (> (len args) 1)) (pop args) (new-sym))
       n (new-sym))
       
  '(let (%i 0 %n %(if v? (pop args) args))
     (while (< %i %n)
       %body..
       (inc %i)))))

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