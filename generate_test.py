"""
File format
    CSV of time, key, val
"""
from hashlib import sha1
from random import randint
from random import random
from random import sample

import numpy as np


UNIVERSE_SIZE = 32000000

def generate_insert(n):
    keys = set()
    with open('tests/insert.csv', 'w') as f:
        for key in sample(xrange(UNIVERSE_SIZE), n):
            val = sha1(str(key)).hexdigest()
            f.write('INSERT,%d,%s%f\n' % (key, val, random()))
            keys.add(key)
    return keys


def generate_steady_state(n, keys):
    start_size = float(len(keys))
    with open('tests/steady.csv', 'w') as f:
        if random()-0.5 > (start_size - len(keys)) / start_size:
            # insert
            key = randint(0, UNIVERSE_SIZE)
            val = sha1(str(key)).hexdigest()
            f.write('INSERT,%d,%s%f\n' % (key, val, random()))
            keys.add(key)
        else:
            # delete
            f.write('DELETE,%d,\n' % (keys.pop(), ))



def generate_poisson(n, lam=1.):
    last_time = 0
    with open('tests/poisson.csv', 'w') as f:
        for _ in xrange(n):
            last_time += np.random.poisson(lam)
            key = randint(0, 1000)
            val = sha1(str(key)).hexdigest()
            f.write('%d,%d,%s%f\n' % (last_time, key, val, random()))


if __name__ == '__main__':
    keys = generate_insert(2 * 10 ** 6)
    generate_steady_state(10 ** 6, keys)

