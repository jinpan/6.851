package main

import (
    "container/list"
    "math"
    "math/rand"
    "time"
)

type Datum struct {
    key int
    val string
}

func Extend(slice []time.Duration, element time.Duration) []time.Duration {
    n := len(slice)
    if n == cap(slice) {
        // Slice is full; must grow.
        // We double its size and add 1, so if the size is zero we still grow.
        newSlice := make([]time.Duration, len(slice), 2*len(slice)+1)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0 : n+1]
    slice[n] = element
    return slice
}


/*
    Simple primality checking function.  Checks to see if any of the integers
    from 2 to sqrt(n) divide n.
*/
func isPrime(n int) bool {
    for i := int(2); i <= int(math.Sqrt(float64(n))); i++ {
        if n % i == 0 {
            return false
        }
    }
    return true
}


/*
    Simple prime generation function.  Takes in a lower and upper bound,
    starts with a random number in that range, checks higher and higher
    numbers until it gets a prime.  If it gets to the upper bound first,
    it tries again with a different starting point.
*/
func getPrime(lower int, upper int) int {
    for start := lower + int(rand.Intn(int(upper - lower))); start < upper; start++ {
        if isPrime(start) {
            return start
        }
    }
    return getPrime(lower, upper)
}


/*
    Compute the potential.
    Locking may severely degrade performance and is avoided. Consequently,
    this his may not be exact because of race conditions, but should be
    numerically stable enough for race conditions to not be a huge deal.
*/

func calcPotential(data []*list.List, n int, m int) float64 {
    potential := 0.0
    expected_length := float64(n) / float64(m)
    cutoff := expected_length + 1.0
    for _, datum := range data {
        if float64(datum.Len()) > cutoff {
            potential += float64(datum.Len()) - cutoff
        }
    }

    return potential
}

