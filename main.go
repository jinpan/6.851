package main

import (
    "fmt"
)

func main() {
    m := uint(1024)
    k := uint(123)
    v := "asdf"

    sht := makeSimpleHashTable(m)
    bht := makeBackgroundHashTable(m)
    fmt.Println("Finished initializing")

    sht.insert(k, v)
    bht.insert(k, v)

    fmt.Println(*sht.get(k))
    fmt.Println(*bht.get(k))
    fmt.Println(*sht.delete(k))
    fmt.Println(*bht.delete(k))
    fmt.Println(sht.get(k))
    fmt.Println(bht.get(k))
    fmt.Println(sht.delete(k))
    fmt.Println(bht.delete(k))
}

