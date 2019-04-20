(load "iter.gf")

(let switch (mac (val alts..)
  (fold (reverse alts)
        (fun (acc alt)
          (let c (head alt))
          
          '(if %(if (= val '_) c '(%(head c) %val %(tail c)..))
             (do %(tail alt))
             %acc))
        _)))