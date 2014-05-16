package main

import (
    "container/list"
    "fmt"
    "math/rand"
)

type SimpleHashTable struct {
    n int  // total size
    p int  // prime number
    a int  // coefficient of the hash function
    m int  // primary table size
    r bool  // whether we automatically rebalance
    data []*SimpleHashTable2  // array of secondary tables
}

type SimpleHashTable2 struct {
    n int  // total size
    p int  // prime number
    a int  // coefficient of the hash function
    m int  // secondary table size
    r bool  // whether we automatically rebalance
    data []*list.List  // actual data
}

func makeSimpleHashTable(m int, rebalance bool) *SimpleHashTable {
    p := getPrime(u, 2*u)
    ht := SimpleHashTable {
        n: 0,
        p: p,
        a: rand.Intn(p),
        m: m,
        r: rebalance,
        data: make([]*SimpleHashTable2, m),
    }

    for i := 0; i < m; i++ {
        p = getPrime(u, 2*u)
        ht.data[i] = &SimpleHashTable2{
            n: 0,
            p: p,
            a: rand.Intn(p),
            m: m,
            r: rebalance,
            data: make([]*list.List, m),
        }

        for j := 0; j < m; j++ {
            ht.data[i].data[j] = list.New()
        }
    }
    return &ht
}

/*
    Hash the key with the params in the table
*/
func (ht *SimpleHashTable) hash(key int) int {
    return ((ht.a * key) % ht.p) % ht.m
}

/*
    Hash the key with the params in the table
*/
func (ht2 *SimpleHashTable2) hash(key int) int {
    return ((ht2.a * key) % ht2.p) % ht2.m
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket
*/
func (ht *SimpleHashTable) insert(key int, val string) {
    datum := Datum{key: key, val: val}

    if (ht.data[ht.hash(key)].insert(datum)) {
        ht.double()
        ht.n++
    }
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket as a Datum object
*/
func (ht2 *SimpleHashTable2) insert(datum Datum) bool {
    key := datum.key
    val := datum.val

    llist := ht2.data[ht2.hash(key)]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            datum := e.Value.(Datum)
            datum.val = val
            return false
        }
    }
    llist.PushBack(Datum{key: key, val: val})
    ht2.n++

    if ht2.r {
        ht2.rebalance()
    }
    return true
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht *SimpleHashTable) get(key int) *string {
    bucket := ht.data[ht.hash(key)]
    return bucket.get(key)
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht2 *SimpleHashTable2) get(key int) *string {
    llist := ht2.data[ht2.hash(key)]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            result := e.Value.(Datum).val
            return &result
        }
    }
    return nil
}

/*
    Deletes the pointer to the value matching the key from the hash table.
*/
func (ht *SimpleHashTable) del(key int) *string {
    bucket := ht.data[ht.hash(key)]
    result := bucket.del(key)
    if result != nil {
        ht.n--
    }
    return result
}

/*
    Deletes the pointer to the value matching the key from the hash table.
*/
func (ht2 *SimpleHashTable2) del(key int) *string {
    llist := ht2.data[ht2.hash(key)]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            result := llist.Remove(e).(Datum).val
            ht2.n--

            if ht2.r {
                ht2.rebalance()
            }
            return &result
        }
    }
    return nil
}

/*
    Doubles the size of the hash table
*/
func (ht *SimpleHashTable) double() {
    if ht.n > ht.m * ht.m {
        fmt.Println("doubling")

        ht.m *= 2
        new_data := make([]*SimpleHashTable2, ht.m)
        for i := 0; i < ht.m; i++ {
            p := getPrime(u, 2*u)
            new_data[i] = &SimpleHashTable2{
                n: 0,
                p: p,
                a: rand.Intn(p),
                m: ht.m,
                r: ht.r,
                data: make([]*list.List, ht.m),
            }
            for j := 0; j < ht.m; j++ {
                new_data[i].data[j] = list.New()
            }
        }

        for i := 0; i < ht.m/2; i++ {  // /2 because old m
            for j := 0; j < ht.data[i].m; j++ {
                llist := ht.data[i].data[j]
                for e := llist.Front(); e != nil; e = e.Next() {
                    datum := e.Value.(Datum)
                    new_data[ht.hash(datum.key)].insert(datum)
                }
            }
        }
        ht.data = new_data
    }
}

/*
    Compute the potential
*/
func (ht2 SimpleHashTable2) calcPotential() float64 {
    potential := 0.0
    expected_length := float64(ht2.n) / float64(ht2.m)
    cutoff := expected_length + 1.0
    for _, datum := range ht2.data {
        if float64(datum.Len()) > cutoff {
            potential += float64(datum.Len()) - cutoff
        }
    }
    return potential
}

/*
    Considers rebalancing the hash table and rebalance if necessary
*/
func (ht2 *SimpleHashTable2) rebalance() {
    for ; ht2.calcPotential() > 19.143 + 0.104 * float64(ht2.n); {
        fmt.Println("rebalancing second level")

        ht2.p = getPrime(u, 2*u)
        ht2.a = rand.Intn(ht2.p)
        data := make([]*list.List, ht2.m)
        for j := 0; j < ht2.m; j++ {
            data[j] = list.New()
        }
        for j := 0; j < ht2.m; j++ {
            for e := ht2.data[j].Front(); e != nil; e = e.Next() {
                datum := e.Value.(Datum)
                data[ht2.hash(datum.key)].PushBack(datum)
            }
        }
        ht2.data = data
    }
}

