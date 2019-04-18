from bench import bench

print(bench(10, '', '''
for _ in range(100000):
  pass
'''))

