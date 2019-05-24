from bench import bench

print(bench(10, '''
class Counter():
    def __init__(self):
        self.n = 0

    def inc(self, d = 1):
        self.n += d
''', '''
for _ in range(10000):
    c = Counter()
    c.inc()
    c.inc(-1)
'''))
