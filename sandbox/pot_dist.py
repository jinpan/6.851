"""
Empirically calculate the potential distribution for different load
factors.  This code will likely be released under the MIT license so
we're pretty okay with this being used for pot distribution, but less
okay with getting sued over pot distribution use.
"""
from random import randint
from scipy.stats import linregress
from scipy.stats.mstats import moment
import numpy as np
import matplotlib.pyplot as plt

TABLE_SIZE = 1024
UNIVERSE_SIZE = 32000000

def compute_potential(data, n):
    expected_collisions = n / float(TABLE_SIZE)
    cutoff = expected_collisions + 1

    potential = 0
    for datum in data:
        potential += max(0, datum - cutoff)
    return potential


def get_potential(alpha, unique_inserts=True):

    num_elements = int(alpha * TABLE_SIZE)
    counts = [0] * TABLE_SIZE

    keys = set()
    for _ in xrange(num_elements):
        if unique_inserts:
            counts[randint(0, TABLE_SIZE-1)] += 1
        else:
            key = randint(0, UNIVERSE_SIZE / TABLE_SIZE)
            if key not in keys:
                counts[randint(0, TABLE_SIZE-1)] += 1
                keys.add(key)
    return compute_potential(counts, num_elements)


if __name__ == '__main__':

    alphas = []
    moments = [[], [], [], []]
    for alpha in (num/10. for num in xrange(1, 51)):
        potentials = [get_potential(alpha) for _ in range(100)]

        print alpha
        alphas.append(alpha)
        moments[0].append(np.mean(potentials))
        moments[1].append(np.std(potentials))
        moments[2].append(moment(potentials, 3))
        moments[3].append(moment(potentials, 4))

    plt.figure()
    plt.xlabel('Load Factor')
    plt.ylabel('Potential')
    plt.title('Potential as a function of Load Factor')
    plt.plot(alphas, moments[0])
    plt.savefig('mean_100_50.png')

    plt.figure()
    plt.plot(moments[1])
    plt.savefig('std_100_50.png')

    slope, intercept, r_val, p_val, std_err = linregress(alphas, moments[0])
    print slope, intercept, r_val, p_val, std_err

    slope, intercept, r_val, p_val, std_err = linregress(alphas, moments[1])
    print slope, intercept, r_val, p_val, std_err

