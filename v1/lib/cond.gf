(load "iter.gf")

(let switch (mac (val alts..)
  (let v? (not (= val '_)))
  
  (fold (reverse alts)
        (fun (acc alt)
          (let c (head alt) v? (and v? (not (= c 'T))))
          
          '(if %(if v? '(%(head c) %val %(tail c)..) c)
             (do %(tail alt))
             %acc))
        _)))