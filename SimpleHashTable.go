package main

import (
    "container/list"
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
        p: getPrime(100 * m, 200 * m),
        m: m,
        data: make([]*SimpleHashTable2, m),
    }

    for i := uint(0); i < m; i++ {
        ht.data[i] = &SimpleHashTable2{
            n: 0,
            p: getPrime(100 * m, 200 * m),
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
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket as a Datum object
*/
func (ht2 *SimpleHashTable2) insert(key uint, val string) {
    llist := ht2.data[(key * ht2.p) % ht2.m]
    llist.PushBack(Datum{key: key, val: val})
    ht2.n += 1
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
            return &result
        }
    }
    return nil
}

