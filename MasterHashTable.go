package main

import (
    "container/list"
    "math/rand"
    "sync"
)

type MasterHashTable struct {
    n int  // total size
    p int  // prime number
    a int  // coefficient of the hash function
    m int  // primary table size
    changed_idx chan int  //  index of changed secondary table
    data []*MasterHashTable2  // array of secondary tables
}

type MasterHashTable2 struct {
    n int  // total size
    p int  // prime number
    a int  // coefficient of the hash function
    m int  // secondary table size
    lock sync.Mutex  // lock
    data []*list.List  // actual data
}

func makeMasterHashTable(m int) *MasterHashTable {
    p := getPrime(u, 2*u)
    ht := MasterHashTable {
        n: 0,
        p: p,
        a: rand.Intn(p),
        m: m,
        changed_idx: make(chan int, 100),
        data: make([]*MasterHashTable2, m),
    }

    for i := 0; i < m; i++ {
        p = getPrime(u, 2*u)
        ht.data[i] = &MasterHashTable2{
            n: 0,
            p: p,
            a: rand.Intn(p),
            m: m,
            lock: sync.Mutex{},
            data: make([]*list.List, m),
        }

        for j := 0; j < m; j++ {
            ht.data[i].data[j] = list.New()
        }
    }
    go ht.run()
    return &ht
}

/*
    Hash the key with the params in the table
*/
func (ht *MasterHashTable) hash(key int) int {
    return ((ht.a * key) % ht.p) % ht.m
}

/*
    Hash the key with the params in the table
*/
func (ht2 *MasterHashTable2) hash(key int) int {
    return ((ht2.a * key) % ht2.p) % ht2.m
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket.
*/
func (ht *MasterHashTable) insert(key int, val string) {
    datum := Datum{key: key, val: val}

    idx := ht.hash(key)
    ht.data[idx].insert(datum)
    ht.n++
    ht.changed_idx <- idx

    ht.double()
}

/*
    Inserts the key/val pair into the hash table.  Gets the appropriate bucket
    and inserts the k/v pair into the bucket as a Datum object.
*/
func (ht2 *MasterHashTable2) insert(datum Datum) {
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
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht *MasterHashTable) get(key int) *string {

    bucket := ht.data[(key * ht.p) % ht.m]
    result := bucket.get(key)

    return result
}

/*
    Retrieves the pointer to the value matching the key from the hash table.
*/
func (ht2 *MasterHashTable2) get(key int) *string {
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
func (ht *MasterHashTable) del(key int) *string {

    idx := ht.hash(key)
    bucket := ht.data[idx]

    result := bucket.del(key)
    if result != nil {
        ht.n--
        ht.changed_idx <- idx
    }

    return result
}

/*
    Deletes the pointer to the value matching the key from the hash table.
*/
func (ht2 *MasterHashTable2) del(key int) *string {
    ht2.lock.Lock()
    defer ht2.lock.Unlock()

    llist := ht2.data[(key * ht2.p) % ht2.m]
    for e := llist.Front(); e != nil; e = e.Next() {
        if e.Value.(Datum).key == key {
            result := llist.Remove(e).(Datum).val
            ht2.n--

            return &result
        }
    }
    return nil
}


func (ht *MasterHashTable) double() {
    if ht.n > ht.m * ht.m {

        ht.m *= 2
        new_data := make([]*MasterHashTable2, ht.m)
        for i := 0; i < ht.m; i++ {
            p := getPrime(u, 2*u)
            new_data[i] = &MasterHashTable2{
                n: 0,
                p: p,
                a: rand.Intn(p),
                m: ht.m,
                lock: sync.Mutex{},
                data: make([]*list.List, ht.m),
            }
            for j := 0; j < ht.m; j++ {
                new_data[i].data[j] = list.New()
            }
        }

        for i := 0; i < ht.m/2; i++ {  // /2 because old m
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
func (ht *MasterHashTable) run() {
    actives := make([]bool, 2048)
    changes := make([]bool, 2048)

    for {
        select {
        case idx := <-ht.changed_idx: {
            if actives[idx] {
                changes[idx] = true
            } else {
                ht2 := ht.data[idx]
                active := actives[idx]
                change := changes[idx]

                active = true
                change = false
                go ht2.cleanup(&active, &change, idx, true)
            }
        }
        /*
        case <-ht.stop: {
            actives = make(map[int] bool)
            changes = make(map[int] bool)
        }
        */
        }
    }
}


func (ht2 *MasterHashTable2) cleanup(active, change *bool, idx int, original bool) {
    potential := calcPotential(ht2.data, ht2.n, ht2.m)
    if potential < 19.143 + 0.104 * float64(ht2.n) {
        *active = false
        return
    }

    for {
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
        if *change { continue }
        new_potential := calcPotential(data, ht2.n, ht2.m)
        if *change { continue }

        if new_potential < 19.143 + 0.104 * float64(ht2.n) {
            ht2.lock.Lock()

            if *change {
                *change = false
                ht2.lock.Unlock()
                ht2.cleanup(active, change, idx, false)
                return
            } else {
                ht2.p = p
                ht2.a = a
                ht2.data = data
                ht2.lock.Unlock()
                *active = false
                return
            }
        }
    }
}

