from sys import argv

import matplotlib.pyplot as plt
import numpy as np


def to_ns(string):
    if string == '0':
        return 0
    if string[-2:] == 'ns':
        return float(string[:-2])
    if string[-2:] == 'us':
        return float(string[:-2]) * 10 ** 3
    if string[-2:] == 'ms':
        return float(string[:-2]) * 10 ** 6
    if string[-1:] == 's':
        return float(string[:-2]) * 10 ** 9

    raise Exception


if __name__ == '__main__':
    ht_type, num_cores = argv[1], argv[2]

    with open('%s_get%s.csv' % (ht_type, num_cores)) as f:
        durations = [to_ns(line.rstrip()) for line in f.readlines()]

    plt.figure()
    plt.hist(durations, range=(0, 1000))
    plt.title("Insert time distribution")
    plt.xlabel("Time (ns)")
    plt.ylabel("Frequency")

    plt.axvline(np.percentile(durations, 50), c='r')
    plt.axvline(np.percentile(durations, 90), c='r')
    plt.axvline(np.percentile(durations, 95), c='r')
    plt.savefig('%s%s_get.png' % (ht_type, num_cores))

