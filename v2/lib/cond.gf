(mac switch (alts..)
  (tr (reverse alts) ()
      (fun (acc alt)
        '(if %(head alt)
           (do %(tail alt)..)
           %acc))))