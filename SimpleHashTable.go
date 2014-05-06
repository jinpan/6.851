package main

import (
    "container/list"
    "math"
    "fmt"
)

type SimpleHashTable struct {
    n uint  // total size
    p uint  // prime number
    m uint  // primary table size
    data []*SimpleHashTable2  // array of secondary tables
}

type SimpleHashTable2 struct {
    n uint  // total size
    p uint  // prime number
    m uint  // secondary table size
    data []*list.List  // actual data
}

func makeSimpleHashTable(m uint) *SimpleHashTable {
    ht := SimpleHashTable {
        n: 0,
        p: getRPrime(m),
        m: m,
        data: make([]*SimpleHashTable2, m),
    }

    for i := uint(0); i < m; i++ {
        ht.data[i] = &SimpleHashTable2{
            n: 0,
            p: getRPrime(m),
            m: m,
            data: make([]*list.List, m),
        }

        for j := uint(0); j < m; j++ {
            ht.data[i].data[j] = list.New()
        }
    }
    return &ht
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket
*/
func (ht *SimpleHashTable) insert(key uint, val string) {
    bucket := ht.data[(key * ht.p) % ht.m]
    bucket.insert(key, val)
    ht.n += 1

    ht.rebalance()
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket as a Datum object
*/
func (ht2 *SimpleHashTable2) insert(key uint, val string) {
    llist := ht2.data[(key * ht2.p) % ht2.m]
    llist.PushBack(Datum{key: key, val: val})
    ht2.n += 1

    ht2.rebalance()
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht *SimpleHashTable) get(key uint) *string {
    bucket := ht.data[(key * ht.p) % ht.m]
    return bucket.get(key)
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht2 *SimpleHashTable2) get(key uint) *string {
    llist := ht2.data[(key * ht2.p) % ht2.m]
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
func (ht *SimpleHashTable) delete(key uint) *string {
    bucket := ht.data[(key * ht.p) % ht.m]
    result := bucket.delete(key)
    if result != nil {
        ht.n -= 1
    }

    ht.rebalance()
    return result
}

/*
    Deletes the pointer to the value matching the key from the hash table.
*/
func (ht2 *SimpleHashTable2) delete(key uint) *string {
    llist := ht2.data[(key * ht2.p) % ht2.m]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            result := llist.Remove(e).(Datum).val
            ht2.n -= 1

            ht2.rebalance()
            return &result
        }
    }
    ht2.rebalance()
    return nil
}

/*
    Compute the potential
*/
func (ht SimpleHashTable) calcPotential() float64 {
    potential := 0.0
    expected_length := float64(ht.n) / float64(ht.m)
    for _, datum := range ht.data {
        if float64(datum.n) > expected_length + 1.0 {
            potential += float64(datum.n) - (expected_length + 1.0)
        }
    }
    return potential
}

/*
    Compute the potential
*/
func (ht2 SimpleHashTable2) calcPotential() float64 {
    potential := 0.0
    expected_length := float64(ht2.n) / float64(ht2.m)
    for _, datum := range ht2.data {
        potential += math.Max(0.0, float64(datum.Len()) - math.Ceil(expected_length))
    }
    return potential
}

/*
    Considers rebalancing the hash table and rebalance if necessary
*/
func (ht *SimpleHashTable) rebalance() {
    for ; ht.calcPotential() > 100; {
        fmt.Println("rebalancing first level", ht.calcPotential(), ht.n)
        data := make([]uint, ht.m)
        for i, datum := range ht.data {
            data[i] = datum.n
        }
        fmt.Println(data)
        new_ht := makeSimpleHashTable(ht.m)
        for i := uint(0); i < ht.m; i++ {
            for j := uint(0); j < ht.data[i].m; j++ {
                llist := ht.data[i].data[j]
                for e := llist.Front(); e != nil; e = e.Next() {
                    datum := e.Value.(Datum)
                    new_ht.insert(datum.key, datum.val)
                }
            }
        }
        ht = new_ht
    }
}


/*
    Considers rebalancing the hash table and rebalance if necessary
*/
func (ht2 *SimpleHashTable2) rebalance() {
    for ; ht2.calcPotential() > 10; {
        fmt.Println("rebalancing second level")
        p := getRPrime(ht2.m)
        data := make([]*list.List, ht2.m)
        for j := uint(0); j < ht2.m; j++ {
            data[j] = list.New()
        }
        for j := uint(0); j < ht2.m; j++ {
            for e := ht2.data[j].Front(); e != nil; e = e.Next() {
                datum := e.Value.(Datum)
                data[(datum.key * p) % ht2.m].PushBack(datum)
            }
        }
        ht2.p = p
        ht2.data = data
    }
}
