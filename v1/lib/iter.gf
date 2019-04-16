(let loop (macro (body..)
  (let done (g-sym) result (g-sym))
  
  '(let (break (macro (args..) '(recall T %args..)))
     ((fun ((%done F) %result..)
        (if %done %result.. (do %body.. (recall))))))))

(let while (macro (cond body..)
  '(loop
     (if %cond _ (break))
     %body..)))