from bench import bench

print(bench(10, '''
def fib(n):
  return n if n < 2 else fib(n-1) + fib(n-2)
''', '''
for _ in range(10):
  fib(20)
'''))
