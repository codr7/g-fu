(let loop (macro (body..)
  (let done? (g-sym) result (g-sym))
  
  '(let (break (macro (args..) '(recall T %args..)))
     ((fun ((%done? F) %result..)
        (if %done? %result.. (do %body.. (recall))))))))

(let while (macro (cond body..)
  '(loop
     (if %cond _ (break))
     %body..)))

(let g-for (macro (arg body..)
  (let i (g-sym) n (g-sym))
  '(let (%i 0 %n %arg)
     (while (< %i %n)
       %body..
       (inc %i)))))