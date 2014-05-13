"""
File format
    CSV of time, key, val
"""
from hashlib import sha1
from random import randint
from random import random

import numpy as np


def generate_simple(n):
    with open('tests/simple.csv', 'w') as f:
        for _ in xrange(n):
            key = randint(0, 32000000)
            val = sha1(str(key)).hexdigest()
            f.write('%d, %s%f\n' % (key, val, random()))


def generate_poisson(n, lam=1.):
    last_time = 0
    with open('tests/poisson.csv', 'w') as f:
        for _ in xrange(n):
            last_time += np.random.poisson(lam)
            key = randint(0, 1000)
            val = sha1(str(key)).hexdigest()
            f.write('%d, %d, %s%f\n' % (last_time, key, val, random()))


if __name__ == '__main__':
    generate_simple(10 ** 6)
    # generate_poisson(10**6, 25.)

