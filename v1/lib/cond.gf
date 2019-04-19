(let g-if (macro (cond x (y _))
  '(or (and %cond %x) %y)))

(let switch (macro (cond alts..)
     'foo
  ))