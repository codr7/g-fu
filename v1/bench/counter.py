from bench import bench

print(bench(10, '''
class Counter():
    def __init__(self):
        self.n = 0

    def inc(self):
        self.n += 1

    def dec(self):
        self.n -= 1
''', '''
for _ in range(1000):
    c = Counter()
    c.inc()
    c.dec()
'''))
