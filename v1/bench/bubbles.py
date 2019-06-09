from bench import bench

print(bench(10, '''
from random import randint

def bubbles(vs):
  done, n = False, len(vs)
  
  while not done:
    done, n = True, n-1
    
    for i in range(n-1):
      x, y = vs[i], vs[i+1]

      if x > y:
        vs[i], vs[i+1] = y, x
        done = False

  return vs 

vals = []

for _ in range(100):
  vals.append(randint(0, 100000))
''', '''
bubbles(vals[:])
'''))
