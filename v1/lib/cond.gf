(let switch (mac (alts..)
  (fold (reverse alts) _
        (fun (acc alt)
          '(if %(head alt)
             (do %(tail alt)..)
             %acc)))))