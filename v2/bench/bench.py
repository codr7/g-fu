from timeit import Timer

def bench(reps, setup, test):
    Timer(test, setup).timeit(reps)
    return int(Timer(test, setup).timeit(reps) * 1000)
