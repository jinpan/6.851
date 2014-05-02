package main

import (
    "container/list"
    "math"
    "sync"
)

type BackgroundHashTable struct {
    n uint  // total size
    p uint  // prime number
    m uint  // primary table size
    data []*BackgroundHashTable2  // array of secondary tables
}

type BackgroundHashTable2 struct {
    n uint  // total size
    p uint  // prime number
    m uint  // secondary table size
    accessed chan bool  // accessed boolean
    lock sync.Mutex  // lock
    data []*list.List  // actual data
}

func makeBackgroundHashTable(m uint) *BackgroundHashTable {
    ht := BackgroundHashTable {
        n: 0,
        p: getPrime(100 * m, 200 * m),
        m: m,
        data : make([]*BackgroundHashTable2, m),
    }

    for i := uint(0); i < m; i++ {
        ht.data[i] = &BackgroundHashTable2{
            n: 0,
            p: getPrime(100 * m, 200 * m),
            m: m,
            accessed: make(chan bool),
            lock: sync.Mutex{},
            data: make([]*list.List, m),
        }

        for j := uint(0); j < m; j++ {
            ht.data[i].data[j] = list.New()
        }
        go cleanup(ht.data[i])
    }
    return &ht
}

func cleanup(ht2 *BackgroundHashTable2) {
    for {
        <-ht2.accessed

        // TODO: optimize this
        for {
            potential := ht2.calcPotential()
            if potential < 10 {  // TODO: Make a better threshold
                break
            }
            // otherwise, try to reduce potential
            ht2.lock.Lock()

            p := getPrime(100 * ht2.m, 200 * ht2.m)
            data := make([]*list.List, ht2.m)
            for j:= uint(0); j < ht2.m; j++ {
                data[j] = list.New()
            }
            for j:= uint(0); j < ht2.m; j++ {
                for e := ht2.data[j].Front(); e != nil; e = e.Next() {
                    datum := e.Value.(Datum)
                    ht2.data[(datum.key * p) % ht2.m].PushBack(datum)
                }
            }
            ht2.data = data
            ht2.lock.Unlock()
        }
    }
}

func (ht2 *BackgroundHashTable2) calcPotential() float64 {
    expected_length := float64(ht2.n) / float64(ht2.m)
    potential := 0.0
    for _, datum := range ht2.data {
        potential += math.Max(0.0, float64(datum.Len()) - expected_length)
    }

    return potential
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket
*/
func (ht *BackgroundHashTable) insert(key uint, val string) {
    bucket := ht.data[(key * ht.p) % ht.m]
    bucket.insert(key, val)
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket as a Datum object
*/
func (ht2 *BackgroundHashTable2) insert(key uint, val string) {
    llist := ht2.data[(key * ht2.p) % ht2.m]
    llist.PushBack(Datum{key: key, val: val})
    ht2.accessed <- true
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht *BackgroundHashTable) get(key uint) *string {
    bucket := ht.data[(key * ht.p) % ht.m]
    return bucket.get(key)
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht2 *BackgroundHashTable2) get(key uint) *string {
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
func (ht *BackgroundHashTable) delete(key uint) *string {
    bucket := ht.data[(key * ht.p) % ht.m]
    return bucket.delete(key)
}

/*
    Deletes the pointer to the value matching the key from the hash table.
*/
func (ht2 *BackgroundHashTable2) delete(key uint) *string {
    llist := ht2.data[(key * ht2.p) % ht2.m]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            result := llist.Remove(e).(Datum).val
            ht2.accessed <- true
            return &result
        }
    }
    return nil
}

