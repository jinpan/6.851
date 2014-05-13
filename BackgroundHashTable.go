package main

import (
    "container/list"
    "fmt"
    "math/rand"
    "sync"
)

type BackgroundHashTable struct {
    n int  // total size
    p int  // prime number
    a int  // coefficient of the hash function
    m int  // primary table size
    data []*BackgroundHashTable2  // array of secondary tables
}

type BackgroundHashTable2 struct {
    n int  // total size
    p int  // prime number
    a int  // coefficient of the hash function
    m int  // secondary table size
    changed chan bool  // changed boolean
    lock sync.Mutex  // lock
    data []*list.List  // actual data
}

func makeBackgroundHashTable(m int) *BackgroundHashTable {
    p := getPrime(u, 2*u)
    ht := BackgroundHashTable {
        n: 0,
        p: p,
        a: rand.Intn(p),
        m: m,
        data : make([]*BackgroundHashTable2, m),
    }

    for i := 0; i < m; i++ {
        p = getPrime(u, 2*u)
        ht.data[i] = &BackgroundHashTable2{
            n: 0,
            p: p,
            a: rand.Intn(p),
            m: m,
            changed: make(chan bool),
            lock: sync.Mutex{},
            data: make([]*list.List, m),
        }

        for j := 0; j < m; j++ {
            ht.data[i].data[j] = list.New()
        }
        go cleanup(ht.data[i])
    }
    return &ht
}

/*
    Hash the key with the params in the table
*/
func (ht *BackgroundHashTable) hash(key int) int {
    return ((ht.a * key) % ht.p) % ht.m
}

/*
    Hash the key with the params in the table
*/
func (ht2 *BackgroundHashTable2) hash(key int) int {
    return ((ht2.a * key) % ht2.p) % ht2.m
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket
*/
func (ht *BackgroundHashTable) insert(key int, val string) {
    datum := Datum{key: key, val: val}

    ht.data[ht.hash(key)].insert(datum)
    ht.n++
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket as a Datum object
*/
func (ht2 *BackgroundHashTable2) insert(datum Datum) {
    key := datum.key
    val := datum.val

    llist := ht2.data[ht2.hash(key)]
    llist.PushBack(Datum{key: key, val: val})
    ht2.n++

    ht2.changed <- true
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht *BackgroundHashTable) get(key int) *string {
    bucket := ht.data[(key * ht.p) % ht.m]
    return bucket.get(key)
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht2 *BackgroundHashTable2) get(key int) *string {
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
func (ht *BackgroundHashTable) del(key int) *string {
    bucket := ht.data[(key * ht.p) % ht.m]
    result := bucket.del(key)
    if result != nil {
        ht.n--
    }
    return result
}

/*
    Deletes the pointer to the value matching the key from the hash table.
*/
func (ht2 *BackgroundHashTable2) del(key int) *string {
    llist := ht2.data[(key * ht2.p) % ht2.m]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            result := llist.Remove(e).(Datum).val
            ht2.n--

            ht2.changed <- true
            return &result
        }
    }
    return nil
}

/*
    Compute the potential
*/
func (ht2 *BackgroundHashTable2) calcPotential() float64 {
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

func cleanup(ht2 *BackgroundHashTable2) {
    for {
        changed := <-ht2.changed
        if !changed {
            return
        }

        for ; ht2.calcPotential() > 19.143 + 0.104 * float64(ht2.n); {
            ht2.lock.Lock()
            fmt.Println("rebalancing second level")
            fmt.Println(ht2.a, ht2.p, ht2.m, ht2.n, ht2.calcPotential())

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
            ht2.lock.Unlock()
        }
    }
}

