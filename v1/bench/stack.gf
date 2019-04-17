(dump, bench 10
  (let, s _,
    for (i 100000) (push s i),
    for 100000 (pop s)))