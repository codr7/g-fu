from bench import bench

print(bench(10, '', '''
s = []

for i in range(100000):
  s.append(i)

for _ in range(100000):
  s.pop()
'''))

