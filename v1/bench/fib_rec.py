from bench import bench

print(bench(10, '''
def fib_rec(n):
  return n if n < 2 else fib_rec(n-1) + fib_rec(n-2)
''', '''
for _ in range(10):
  fib_rec(20)
'''))
