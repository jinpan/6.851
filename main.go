package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "time"
)

var (
    u = 32000000
)

func main() {

    ht := makeSimpleHashTable(16, false)
    // ht := makeBackgroundHashTable(16)
    // ht := makeMasterHashTable(16)


    // Insertion
    insert_file, _ := os.Open("tests/insert.csv")
    insert_reader := csv.NewReader(insert_file)
    insert_times := make([]time.Duration, 0)

    for i := 0; ; i++ {
        data, err := insert_reader.Read()
        if err != nil {
            break
        }

        key, _ := strconv.ParseInt(data[1], 0, 0)
        val := data[2]
        start := time.Now()
        ht.insert(int(key), val)
        insert_times = Extend(insert_times, time.Now().Sub(start))
    }

    insert_out, _ := os.Create("data/simple4.csv")
    for _, datum := range insert_times {
        insert_out.Write([]byte(datum.String() + "\n"))
    }

    // Get
    get_file, _ := os.Open("tests/get.csv")
    get_reader := csv.NewReader(get_file)
    get_times := make([]time.Duration, 0)
    get_start := time.Now()

    for i := 0; ; i++ {
        data, err := get_reader.Read()
        if err != nil {
            break
        }

        key, _ := strconv.ParseInt(data[1], 0, 0)
        start := time.Now()
        ht.get(int(key))
        get_times = Extend(get_times, time.Now().Sub(start))
    }
    fmt.Println("Elapsed time", time.Now().Sub(get_start))

    get_out, _ := os.Create("data/simple_get4.csv")
    for _, datum := range get_times {
        get_out.Write([]byte(datum.String() + "\n"))
    }


}

