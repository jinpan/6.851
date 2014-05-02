package main

import (
    "math"
    "math/rand"
)

type Datum struct {
    key uint
    val string
}

/*
    Simple primality checking function.  Checks to see if any of the integers
    from 2 to sqrt(n) divide n.
*/
func isPrime(n uint) bool {
    for i := uint(2); i <= uint(math.Sqrt(float64(n))); i++ {
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
func getPrime(lower uint, upper uint) uint {
    for start := lower + uint(rand.Intn(int(upper - lower))); start < upper; start++ {
        if isPrime(start) {
            return start
        }
    }
    return getPrime(lower, upper)
}

