package main

import (
    "container/list"
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
            changed: make(chan bool, 10),
            lock: sync.Mutex{},
            data: make([]*list.List, m),
        }

        for j := 0; j < m; j++ {
            ht.data[i].data[j] = list.New()
        }
        go ht.data[i].cleanup()
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
    ht.double()
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket as a Datum object
*/
func (ht2 *BackgroundHashTable2) insert(datum Datum) {
    key := datum.key
    val := datum.val

    ht2.lock.Lock()
    llist := ht2.data[ht2.hash(key)]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            datum := e.Value.(Datum)
            datum.val = val
            ht2.lock.Unlock()
            return
        }
    }
    llist.PushBack(Datum{key: key, val: val})

    ht2.n++
    ht2.lock.Unlock()

    ht2.changed <- true
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht *BackgroundHashTable) get(key int) *string {

    bucket := ht.data[(key * ht.p) % ht.m]
    result := bucket.get(key)

    return result
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht2 *BackgroundHashTable2) get(key int) *string {
    ht2.lock.Lock()
    defer ht2.lock.Unlock()

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
    ht2.lock.Lock()
    defer ht2.lock.Unlock()

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


func (ht *BackgroundHashTable) double() {
    if ht.n > ht.m * ht.m {

        ht.m *= 2
        new_data := make([]*BackgroundHashTable2, ht.m)
        for i := 0; i < ht.m; i++ {
            p := getPrime(u, 2*u)
            new_data[i] = &BackgroundHashTable2{
                n: 0,
                p: p,
                a: rand.Intn(p),
                m: ht.m,
                changed: make(chan bool, 10),
                lock: sync.Mutex{},
                data: make([]*list.List, ht.m),
            }
            for j := 0; j < ht.m; j++ {
                new_data[i].data[j] = list.New()
            }
            go new_data[i].cleanup()
        }

        for i := 0; i < ht.m/2; i++ {  // /2 because old m
            ht.data[i].changed <- false
            for j := 0; j < ht.data[i].m; j++ {
                ht.data[i].lock.Lock()
                llist := ht.data[i].data[j]
                for e := llist.Front(); e != nil; e = e.Next() {
                    datum := e.Value.(Datum)
                    new_data[ht.hash(datum.key)].insert(datum)
                }
                ht.data[i].lock.Unlock()
            }
        }
        ht.data = new_data
    }
}

/*
    Reduces the potential
*/
func (ht2 *BackgroundHashTable2) cleanup() {
    for {
        if !<-ht2.changed {
            return
        }
        for {
            potential := calcPotential(ht2.data, ht2.n, ht2.m)
            if potential < 19.143 + 0.104 * float64(ht2.n) {
                break
            }

            p := getPrime(u, 2*u)
            a := rand.Intn(p)
            data := make([]*list.List, ht2.m)
            for j := 0; j < ht2.m; j++ {
                data[j] = list.New()
            }
            for j := 0; j < ht2.m; j++ {
                for e := ht2.data[j].Front(); e != nil; e = e.Next() {
                    datum := e.Value.(Datum)
                    data[((a * datum.key) % p) % ht2.m].PushBack(datum)
                }
            }
            new_potential := calcPotential(data, ht2.n, ht2.m)
            if new_potential > 19.143 + 0.104 * float64(ht2.n) {
                select {
                case cont := <-ht2.changed: {
                    if cont {
                        continue
                    } else {
                        return
                    }
                }
                default: {
                    continue
                }
                }
            } else {
                select {
                case cont := <-ht2.changed: {
                    if cont {
                        continue
                    } else {
                        return
                    }
                }
                default: {
                    ht2.lock.Lock()
                    ht2.p = p
                    ht2.a = a
                    ht2.data = data
                    ht2.lock.Unlock()
                }
                }
            }
        }
    }
}

