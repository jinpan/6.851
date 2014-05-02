package main

import (
    "container/list"
    "fmt"
    "time"
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
    expected_length := float64(ht2.n) / float64(ht2.m)
    for {
        // calculate potential
        potential := 0.
        for _, datum := range ht2.data {
            if float64(datum.Len()) > expected_length {
                potential += float64(datum.Len()) - expected_length
            }
        }
        fmt.Println(potential)

        if potential > 10 {  // TODO: make a better threshold
            // do something to reduce potential
        } else {
            // chill
            time.Sleep(1000)
        }
    }

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
func (ht *BackgroundHashTable2) delete(key uint) *string {
    llist := ht.data[(key * ht.p) % ht.m]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            result := llist.Remove(e).(Datum).val
            return &result
        }
    }
    return nil
}

